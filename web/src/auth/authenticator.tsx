import { setInterceptor } from "@/api/api";
import userApi, { User } from "@/api/user";

export interface AuthProvider {
  isAuthenticated: boolean;
  user: null | User;
  iid: number;
  signin(username: string, password: string): Promise<User | null>;
  register: (username: string, password: string) => Promise<User>;
  authorize: () => Promise<User>;
  signout(): Promise<void>;
}

function validateUser(user: User | null): user is User {
  if (user === null) {
    return false;
  }
  if (user.token === undefined || user.token === "") {
    return false;
  }
  if (user.userId === undefined || user.userId <= 0) {
    return false;
  }
  return true;
}

export const Authenticator: AuthProvider = {
  isAuthenticated: false,
  user: null,
  iid: 0,
  async signin(username: string, password: string) {
    try {
      const user = await userApi.login(username, password);
      switch (true) {
        case user === undefined:
        case user.token === undefined:
        case user.token === "":
          console.log("User Sign In Failed", user);
          Authenticator.isAuthenticated = false;
          return null;
        default:
          console.log("User Sign In Success", user);
          Authenticator.isAuthenticated = true;
          Authenticator.user = user;
          console.log("Authenticator.iid", Authenticator.iid);
          if (Authenticator.iid === 0) {
            Authenticator.iid = setInterceptor(user.token);
          }
          break;
      }
      return user;
    } catch (error: any) {
      console.log("User Sign In Failed", error.message);
      return null;
    }
  },
  async signout() {
    await userApi.logout();
    Authenticator.isAuthenticated = false;
    Authenticator.user = null;
  },
  async register(username: string, password: string) {
    const user = await userApi.register(username, password);
    if (user.token !== undefined && user.token !== "") {
      if (Authenticator.iid === 0) {
        Authenticator.iid = setInterceptor(user.token);
      }
      Authenticator.isAuthenticated = true;
    }
    return user;
  },
  async authorize() {
    let user: User = {
      userId: 0,
      token: "",
    };
    try {
      user = await userApi.authorize();
      console.log("authorize response", user);

      if (validateUser(user)) {
        Authenticator.user = user;
        Authenticator.isAuthenticated = true;
        if (Authenticator.iid === 0) {
          Authenticator.iid = setInterceptor(user.token);
        }
      } else {
        Authenticator.user = null;
        Authenticator.isAuthenticated = false;
      }
    } catch (error) {
      console.error("Authorization failed:", error);
      Authenticator.isAuthenticated = false;
    } finally {
      return user;
    }
  },
};
