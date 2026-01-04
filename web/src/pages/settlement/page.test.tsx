import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { redirect } from "react-router";
import { describe, it, expect, vi } from "vitest";
import { SettlementPage, TimelineTab, PopulationTab, StorageTab } from "~/pages/settlement/index";
import { TestRouterWithRoutes } from "~/test/TestRouter";

function renderSettlementPage(initialPath = "/settlements/1") {
  render(
    <TestRouterWithRoutes
      initialEntries={[initialPath]}
      routes={[
        { path: "/", element: <div>Home</div> },
        { path: "/settlements", element: <div>Settlement Select</div> },
        {
          path: "/settlements/:settlementId",
          element: <SettlementPage />,
          children: [
            {
              index: true,
              loader: ({ params }) => redirect(`/settlements/${params.settlementId}/timeline`),
            },
            { path: "timeline", element: <TimelineTab /> },
            { path: "population", element: <PopulationTab />, loader: () => [] },
            { path: "storage", element: <StorageTab /> },
          ],
        },
      ]}
    />
  );
  return userEvent.setup();
}

describe("SettlementPage", () => {
  describe("navigation", () => {
    it("defaults to the timeline tab", async () => {
      renderSettlementPage("/settlements/1");

      await screen.findByRole("button", { name: /timeline/i });
      expect(document.getElementById("timeline")).toBeInTheDocument();
    });

    it("can navigate to the storage tab", async () => {
      const user = renderSettlementPage("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /storage/i }));

      expect(document.getElementById("storage")).toBeInTheDocument();
    });

    it("can navigate to the population tab", async () => {
      const user = renderSettlementPage("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /population/i }));

      expect(document.getElementById("population")).toBeInTheDocument();
    });

    it("can navigate to settlement selection", async () => {
      const user = renderSettlementPage("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("link", { name: /settlements/i })
      );

      expect(await screen.findByText("Settlement Select")).toBeInTheDocument();
    });

    it("can log out", async () => {
      const mockFetch = vi.fn().mockResolvedValue({ ok: true });
      vi.stubGlobal("fetch", mockFetch);

      const user = renderSettlementPage("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("button", { name: /logout/i })
      );

      expect(mockFetch).toHaveBeenCalledWith("/api/auth/logout", {
        credentials: "include",
      });
      expect(await screen.findByText("Home")).toBeInTheDocument();

      vi.unstubAllGlobals();
    });
  });
});
