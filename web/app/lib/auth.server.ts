import { redirect } from "react-router";

const API_BASE_URL = process.env.API_URL;
const AUTH_LOGIN_URL = "/auth/login";

export async function checkAuth(request: Request): Promise<boolean> {
  const response = await fetch(`${API_BASE_URL}/auth/check`, {
    headers: {
      cookie: request.headers.get("cookie") || "",
    },
  });
  return response.ok;
}

export async function requireAuth(request: Request): Promise<void> {
  const isAuthenticated = await checkAuth(request);
  if (!isAuthenticated) {
    throw redirect(AUTH_LOGIN_URL);
  }
}
