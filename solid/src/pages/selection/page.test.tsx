import { render, screen, within, waitFor } from "@solidjs/testing-library";
import { MemoryRouter, Route, createMemoryHistory, revalidate } from "@solidjs/router";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, it, expect, beforeEach } from "vitest";
import { Suspense } from "solid-js";
import { server } from "~test/setup";
import SelectSettlement, { getSettlements } from "./page";
import type { SettlementId } from "~/lib/settlement";

const mockSettlements: SettlementId[] = [
  { id: "1", name: "People of the Lantern" },
  { id: "2", name: "The Sun Stalkers" },
  { id: "3", name: "Children of the Stars" },
];

function renderPage() {
  const history = createMemoryHistory();
  render(() => (
    <MemoryRouter history={history}>
      <Route path="/" component={() => (
        <Suspense fallback={<div>Loading...</div>}>
          <SelectSettlement />
        </Suspense>
      )} />
    </MemoryRouter>
  ));
  return userEvent.setup();
}

describe("SelectSettlement", () => {
  beforeEach(() => {
    revalidate(getSettlements.key, true);
  });
  describe("viewing settlements", () => {
    it("shows an empty list when no settlements exist", async () => {
      server.use(http.get("/api/settlements", () => HttpResponse.json([])));

      renderPage();

      await screen.findByRole("button", { name: /create settlement/i });
      expect(screen.queryByRole("link")).not.toBeInTheDocument();
    });

    it("shows all my settlements", async () => {
      server.use(
        http.get("/api/settlements", () => HttpResponse.json(mockSettlements))
      );

      renderPage();

      expect(await screen.findByText("People of the Lantern")).toBeInTheDocument();
      expect(screen.getByText("The Sun Stalkers")).toBeInTheDocument();
      expect(screen.getByText("Children of the Stars")).toBeInTheDocument();
    });
  });

  describe("selecting a settlement", () => {
    it("navigates to the settlement when clicking play", async () => {
      server.use(
        http.get("/api/settlements", () => HttpResponse.json(mockSettlements))
      );

      renderPage();

      await screen.findByText("People of the Lantern");

      const playLinks = screen.getAllByRole("link");
      expect(playLinks[0]).toHaveAttribute("href", "/settlements/1");
      expect(playLinks[1]).toHaveAttribute("href", "/settlements/2");
      expect(playLinks[2]).toHaveAttribute("href", "/settlements/3");
    });
  });

  describe("creating a settlement", () => {
    it("requires a name of at least 5 characters", async () => {
      server.use(http.get("/api/settlements", () => HttpResponse.json([])));

      const user = renderPage();

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
      const settlements: SettlementId[] = [];
      server.use(
        http.get("/api/settlements", () => HttpResponse.json(settlements)),
        http.post("/api/settlements", async ({ request }) => {
          const body = (await request.json()) as { name: string };
          const newSettlement = { id: crypto.randomUUID(), name: body.name };
          settlements.push(newSettlement);
          return HttpResponse.json(newSettlement, { status: 201 });
        })
      );

      const user = renderPage();

      await user.click(
        await screen.findByRole("button", { name: /create settlement/i })
      );

      const dialog = screen.getByRole("dialog");
      await user.type(
        within(dialog).getByPlaceholderText(/enter settlement name/i),
        "New Settlement"
      );
      await user.click(within(dialog).getByRole("button", { name: /^create$/i }));

      await waitFor(() => {
        expect(dialog).not.toHaveAttribute("open");
      });

      expect(await screen.findByText("New Settlement")).toBeInTheDocument();
    });

    it("can be cancelled", async () => {
      server.use(http.get("/api/settlements", () => HttpResponse.json([])));

      const user = renderPage();

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
