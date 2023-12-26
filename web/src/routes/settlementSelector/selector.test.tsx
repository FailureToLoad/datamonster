import Selector from "./selector";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import api, { Settlement } from "@/api/settlement";
import { useLoaderData } from "react-router-dom";
import userEvent from "@testing-library/user-event";

vi.mock("react-router-dom");
vi.mock("@/api/settlement");
const settlement1: Settlement = {
  id: "1",
  name: "First",
  limit: 1,
  departing: 0,
  cc: 0,
  year: 1,
};
const settlement2: Settlement = {
  id: "2",
  name: "Second",
  limit: 1,
  departing: 0,
  cc: 0,
  year: 1,
};
const settlements: Array<Settlement> = [settlement1, settlement2];
const build = () => {
  return <Selector testId="selector" />;
};

describe("Selector", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });
  it("renders a list of settlements", async () => {
    vi.mocked(useLoaderData).mockReturnValue(settlements);
    render(build());
    const selector = screen.getByTestId("selector");
    expect(selector).toBeInTheDocument();
    const listItems = screen.getAllByRole("listitem");
    expect(listItems.length).toEqual(3);
  });
  it("renders the add button if no settlements are loaded", async () => {
    vi.mocked(useLoaderData).mockReturnValue(null);
    render(build());
    const selector = screen.getByTestId("selector");
    expect(selector).toBeInTheDocument();
    const listItems = screen.getAllByRole("listitem");
    expect(listItems.length).toEqual(1);
    const listItem = listItems[0];
    const button = listItem.firstChild;
    expect(button).toHaveAttribute(
      "aria-label",
      "launch add settlement dialogue",
    );
  });
  describe("create settlement dialogue", () => {
    it("launches when the + button is clicked", async () => {
      vi.mocked(useLoaderData).mockReturnValue(null);
      render(build());
      const selector = screen.getByTestId("selector");
      expect(selector).toBeInTheDocument();
      const listItems = screen.getAllByRole("listitem");
      expect(listItems.length).toEqual(1);
      const listItem = listItems[0];
      const button = listItem.firstChild;
      expect(button).toHaveAttribute(
        "aria-label",
        "launch add settlement dialogue",
      );
      await userEvent.click(button as Element);
      const dialogue = screen.getByText("Enter settlement details.");
      expect(dialogue).toBeInTheDocument();
    });
    it("submits a settlement creation request", async () => {
      vi.mocked(useLoaderData).mockReturnValue(null);
      render(build());
      const selector = screen.getByTestId("selector");
      expect(selector).toBeInTheDocument();
      const listItems = screen.getAllByRole("listitem");
      expect(listItems.length).toEqual(1);
      const listItem = listItems[0];
      const button = listItem.firstChild;
      expect(button).toHaveAttribute(
        "aria-label",
        "launch add settlement dialogue",
      );
      await userEvent.click(button as Element);
      const nameField = screen.getByLabelText(/settlement name/i);
      expect(nameField).toBeInTheDocument();
      expect(nameField).toHaveValue("");
      await userEvent.click(nameField);
      await userEvent.paste("First");
      expect(nameField).toHaveValue("First");
      vi.mocked(api.createSettlement).mockResolvedValue(settlement1);
      const addButton = screen.getByRole("button", { name: "Add" });
      expect(addButton).toBeInTheDocument();
      const dialogue = screen.getByRole("dialog");
      expect(dialogue).toBeInTheDocument();
      await userEvent.click(addButton);
      expect(api.createSettlement).toHaveBeenCalledWith({ name: "First" });
      expect(dialogue).not.toBeInTheDocument();
      expect(screen.getByText("First")).toBeInTheDocument();
    });
  });
});
