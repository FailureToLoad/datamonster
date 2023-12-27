import { setInterceptor } from "@/api/api";
import userApi, { User } from "@/api/user";
import { LoaderFunctionArgs, redirect } from "react-router-dom";

export interface AuthProvider {
  authenticated: boolean;
  user: null | User;
  iid: number;
  isAuthenticated(): boolean;
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

export async function AuthLoader({ request }: LoaderFunctionArgs) {
  const user = await Authenticator.authorize();
  if (!Authenticator.authenticated) {
    let params = new URLSearchParams();
    params.set("from", new URL(request.url).pathname);
    return redirect("/login?" + params.toString());
  }
  return user;
}

export const Authenticator: AuthProvider = {
  authenticated: false,
  user: null,
  iid: 0,
  isAuthenticated() {
    return Authenticator.authenticated;
  },
  async signin(username: string, password: string) {
    try {
      const user = await userApi.login(username, password);
      switch (true) {
        case user === undefined:
        case user.token === undefined:
        case user.token === "":
          console.log("User Sign In Failed", user);
          Authenticator.authenticated = false;
          return null;
        default:
          console.log("User Sign In Success", user);
          Authenticator.authenticated = true;
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
    Authenticator.authenticated = false;
    Authenticator.user = null;
  },
  async register(username: string, password: string) {
    const user = await userApi.register(username, password);
    if (user.token !== undefined && user.token !== "") {
      if (Authenticator.iid === 0) {
        Authenticator.iid = setInterceptor(user.token);
      }
      Authenticator.authenticated = true;
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

      if (validateUser(user)) {
        Authenticator.user = user;
        Authenticator.authenticated = true;
        if (Authenticator.iid === 0) {
          Authenticator.iid = setInterceptor(user.token);
        }
      } else {
        Authenticator.user = null;
        Authenticator.authenticated = false;
      }
    } catch (error) {
      console.error("Authorization failed:", error);
      Authenticator.authenticated = false;
    } finally {
      return user;
    }
  },
};
