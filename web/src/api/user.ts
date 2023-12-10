import { User } from "lucide-react";
import api from "./api";
export type User = {
  userId: number;
  token: string;
};

type userApi = {
  login: (username: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
  register: (username: string, password: string) => Promise<User>;
  authorize: () => Promise<User>;
  isLoading: boolean;
};

const requester: userApi = {
  isLoading: false,
  login: async function (username: string, password: string): Promise<User> {
    const response = await api.post<User>("http://localhost:8080/login", {
      username: username,
      password: password,
    });

    return response.data;
  },
  register: async function (username: string, password: string): Promise<User> {
    const response = await api.post<User>("http://localhost:8080/register", {
      username: username,
      password: password,
    });

    return response.data;
  },
  logout: async function (): Promise<void> {
    await api.post("http://localhost:8080/logout");
  },
  authorize: async function (): Promise<User> {
    try {
      const response = await api.post<User>("http://localhost:8080/auth");
      return response.data;
    } catch (error) {
      console.log(`error ${error}`);
      const nullUser: User = { userId: 0, token: "" };
      return nullUser;
    }
  },
};

export default requester;
