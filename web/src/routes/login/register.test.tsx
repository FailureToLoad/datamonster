import { expect, it, describe, vi, beforeEach } from "vitest";
import { RegisterDialogue } from "./register";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import userApi from "@/api/user";
import { setInterceptor } from "@/api/api";
const user = {
  userId: 1,
  token: "token",
};
const mockRegister = vi.fn().mockResolvedValue(user);
vi.mock("@/api/user");
const build = () => {
  return <RegisterDialogue />;
};

describe("Register", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });
  it("renders the register dialogue", async () => {
    render(build());
    const registerButton = screen.getByRole("button", { name: "Register" });
    expect(registerButton).toBeInTheDocument();
    await userEvent.click(registerButton);
    const usernameField = screen.getByLabelText("Username");
    expect(usernameField).toBeInTheDocument();
    const passwordField = screen.getByLabelText("Password");
    expect(passwordField).toBeInTheDocument();
    const confirmPasswordField = screen.getByLabelText(/confirm password/i);
    expect(confirmPasswordField).toBeInTheDocument();
  });
});
