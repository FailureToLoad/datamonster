import { Authenticator } from "@/auth/authenticator";
import login, { LoginAction } from "./login";

export async function LoginLoader() {
  return Authenticator.isAuthenticated();
}
export default login;
export { LoginAction };
