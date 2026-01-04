import { screen, waitFor } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { testApp } from "~test/builder";

describe("SettlementPage", () => {
  describe("navigation", () => {
    it("defaults to the timeline tab", async () => {
      testApp().renderAt("/settlements/1");

      await screen.findByRole("button", { name: /timeline/i });
      await waitFor(() => {
        expect(document.getElementById("timeline")).toBeInTheDocument();
      });
    });

    it("can navigate to the storage tab", async () => {
      const user = testApp().renderAt("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /storage/i }));

      await waitFor(() => {
        expect(document.getElementById("storage")).toBeInTheDocument();
      });
    });

    it("can navigate to the population tab", async () => {
      const user = testApp().renderAt("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /population/i }));

      await waitFor(() => {
        expect(document.getElementById("population")).toBeInTheDocument();
      });
    });

    it("can navigate to settlement selection", async () => {
      const user = testApp().renderAt("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("link", { name: /settlements/i })
      );

      expect(await screen.findByRole("button", { name: /create settlement/i })).toBeInTheDocument();
    });

    it("can log out", async () => {
      const user = testApp().renderAt("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("button", { name: /logout/i })
      );

      expect(await screen.findByText("Home")).toBeInTheDocument();
    });
  });
});
