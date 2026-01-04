import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { SelectSettlement } from "~/pages/selection/page";
import { TestRouter } from "~/test/TestRouter";
import type { SettlementId } from "~/lib/settlement";

const mockSettlements: SettlementId[] = [
  { id: "1", name: "People of the Lantern" },
  { id: "2", name: "The Sun Stalkers" },
  { id: "3", name: "Children of the Stars" },
];

function renderPage(settlements: SettlementId[] = []) {
  render(
    <TestRouter loaderData={settlements}>
      <SelectSettlement />
    </TestRouter>
  );
  return userEvent.setup();
}

describe("SelectSettlement", () => {
  describe("viewing settlements", () => {
    it("shows an empty list when no settlements exist", async () => {
      renderPage([]);

      await screen.findByRole("button", { name: /create settlement/i });
      expect(screen.queryByRole("link")).not.toBeInTheDocument();
    });

    it("shows all my settlements", async () => {
      renderPage(mockSettlements);

      expect(await screen.findByText("People of the Lantern")).toBeInTheDocument();
      expect(screen.getByText("The Sun Stalkers")).toBeInTheDocument();
      expect(screen.getByText("Children of the Stars")).toBeInTheDocument();
    });
  });

  describe("selecting a settlement", () => {
    it("navigates to the settlement when clicking play", async () => {
      renderPage(mockSettlements);

      await screen.findByText("People of the Lantern");

      const playLinks = screen.getAllByRole("link");
      expect(playLinks[0]).toHaveAttribute("href", "/settlements/1");
      expect(playLinks[1]).toHaveAttribute("href", "/settlements/2");
      expect(playLinks[2]).toHaveAttribute("href", "/settlements/3");
    });
  });

  describe("creating a settlement", () => {
    it("requires a name of at least 5 characters", async () => {
      const user = renderPage([]);

      await user.click(
        await screen.findByRole("button", { name: /create settlement/i })
      );

      const dialog = screen.getByRole("dialog");
      const nameInput = within(dialog).getByPlaceholderText(/enter settlement name/i);
      const submitButton = within(dialog).getByRole("button", { name: /^create$/i });

      await user.type(nameInput, "abcd");
      expect(submitButton).toBeDisabled();

      await user.type(nameInput, "e");
      expect(submitButton).toBeEnabled();
    });

    it("creates a settlement and refreshes the list", async () => {
      const mockFetch = vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ id: "new-id", name: "New Settlement" }),
      });
      vi.stubGlobal("fetch", mockFetch);

      const user = renderPage([]);

      await user.click(
        await screen.findByRole("button", { name: /create settlement/i })
      );

      const dialog = screen.getByRole("dialog");
      await user.type(
        within(dialog).getByPlaceholderText(/enter settlement name/i),
        "New Settlement"
      );
      await user.click(within(dialog).getByRole("button", { name: /^create$/i }));

      expect(mockFetch).toHaveBeenCalledWith(
        "/api/settlements",
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({ name: "New Settlement" }),
        })
      );

      vi.unstubAllGlobals();
    });

    it("can be cancelled", async () => {
      const user = renderPage([]);

      await user.click(
        await screen.findByRole("button", { name: /create settlement/i })
      );

      const dialog = screen.getByRole("dialog");
      await user.type(
        within(dialog).getByPlaceholderText(/enter settlement name/i),
        "Draft Settlement"
      );
      await user.click(within(dialog).getByRole("button", { name: /cancel/i }));

      expect(dialog).not.toHaveAttribute("open");
    });
  });
});
