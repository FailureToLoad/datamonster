import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RouterProvider, createMemoryRouter } from "react-router-dom";
import Login, { LoginAction, LoginLoader } from "./login";
import { Authenticator } from "@/auth/authenticator";
import { beforeEach } from "node:test";
import userApi from "@/api/user";

vi.mock("@/api/user");
const user = {
  userId: 1,
  token: "token",
};
const notAuthenticated = vi.fn().mockReturnValue(false);
const mockLogin = vi.fn().mockResolvedValue(user);
const routes = [
  {
    path: "/login",
    element: <Login />,
    loader: LoginLoader,
    action: LoginAction,
  },
];
const router = createMemoryRouter(routes, {
  initialEntries: ["/", "/login"],
  initialIndex: 1,
});
const build = () => {
  return <RouterProvider router={router} />;
};

describe("Login", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });
  it("renders the login page", () => {
    vi.mocked(Authenticator).isAuthenticated = notAuthenticated;
    render(build());
    const usernameField = screen.getByLabelText(/username/i);
    expect(usernameField).toBeInTheDocument();
    const passwordField = screen.getByLabelText(/password/i);
    expect(passwordField).toBeInTheDocument();
    const submitButton = screen.getByRole("button", { name: "Sign In" });
    expect(submitButton).toBeInTheDocument();
  });
  it("submits the login form", async () => {
    vi.mocked(Authenticator).isAuthenticated = notAuthenticated;
    vi.mocked(userApi).login = mockLogin;
    render(build());
    const usernameField = screen.getByLabelText(/username/i);
    expect(usernameField).toHaveValue("");
    await userEvent.click(usernameField);
    await userEvent.paste("username");
    expect(usernameField).toHaveValue("username");

    const passwordField = screen.getByLabelText(/password/i);
    expect(passwordField).toHaveValue("");
    await userEvent.click(passwordField);
    await userEvent.paste("password");
    expect(passwordField).toHaveValue("password");

    const submitButton = screen.getByRole("button", { name: "Sign In" });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);
    expect(mockLogin).toHaveBeenCalledWith("username", "password");
  });
  it("launches the register dialogue", async () => {
    vi.mocked(Authenticator).isAuthenticated = notAuthenticated;
    render(build());
    const registerButton = screen.getByRole("button", {
      name: "Register",
    });
    expect(registerButton).toBeInTheDocument();
    await userEvent.click(registerButton);
    const registerDialogue = screen.getByRole("dialog");
    expect(registerDialogue).toBeInTheDocument();
    const dialogueId = registerDialogue.getAttribute("id");
    expect(dialogueId).toBe("register");
  });
});
