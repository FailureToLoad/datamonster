import { render, screen, waitFor } from "@solidjs/testing-library";
import { MemoryRouter, Route, createMemoryHistory, revalidate } from "@solidjs/router";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, it, expect, beforeEach } from "vitest";
import { Suspense, lazy, createRoot } from "solid-js";
import { server } from "~test/setup";
import { GlossaryContext, type GlossaryData } from "~/hooks/glossary";
import { SettlementPage } from "./page";
import { getSurvivors } from "./population/tab";

const mockGlossary: GlossaryData = {
  disorders: [],
  fightingArts: [],
};

const timeline = lazy(() => import("./timeline/tab"));
const population = lazy(() => import("./population/tab"));
const storage = lazy(() => import("./storage/tab"));

function renderPage(initialPath = "/settlements/1") {
  server.use(
    http.get("/api/settlements/:settlementId/survivors", () => {
      return HttpResponse.json([]);
    })
  );

  render(() => {
    const history = createMemoryHistory();
    history.set({ value: initialPath });
    return (
      <GlossaryContext.Provider value={() => mockGlossary}>
        <MemoryRouter history={history}>
          <Route path="/settlements" component={() => (
            <div>
              <button>Create Settlement</button>
            </div>
          )} />
          <Route path="/settlements/:settlementId" component={(props) => (
            <Suspense fallback={<div>Loading...</div>}>
              <SettlementPage>{props.children}</SettlementPage>
            </Suspense>
          )}>
            <Route path="/timeline" component={timeline} />
            <Route path="/population" component={population} />
            <Route path="/storage" component={storage} />
          </Route>
        </MemoryRouter>
      </GlossaryContext.Provider>
    );
  });
  return userEvent.setup();
}

describe("SettlementPage", () => {
  beforeEach(() => {
    createRoot((dispose) => {
      revalidate(getSurvivors.key, true);
      dispose();
    });
  });

  describe("navigation", () => {
    it("can navigate to the timeline tab", async () => {
      const user = renderPage("/settlements/1/storage");

      const menuButton = await screen.findByRole("button", { name: /storage/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /timeline/i }));

      await waitFor(() => {
        expect(document.getElementById("timeline")).toBeInTheDocument();
      });
    });

    it("can navigate to the storage tab", async () => {
      const user = renderPage("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /storage/i }));

      await waitFor(() => {
        expect(document.getElementById("storage")).toBeInTheDocument();
      });
    });

    it("can navigate to the population tab", async () => {
      const user = renderPage("/settlements/1/timeline");

      const menuButton = await screen.findByRole("button", { name: /timeline/i });
      await user.click(menuButton);

      await user.click(screen.getByRole("link", { name: /population/i }));

      await waitFor(() => {
        expect(document.getElementById("population")).toBeInTheDocument();
      });
    });

    it("can navigate to settlement selection", async () => {
      const user = renderPage("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("link", { name: /settlements/i })
      );

      expect(await screen.findByRole("button", { name: /create settlement/i })).toBeInTheDocument();
    });

    it("can log out", async () => {
      let logoutCalled = false;
      server.use(
        http.get("/auth/logout", () => {
          logoutCalled = true;
          return HttpResponse.json({});
        })
      );

      const user = renderPage("/settlements/1/timeline");

      await user.click(
        await screen.findByRole("button", { name: /logout/i })
      );

      await waitFor(() => {
        expect(logoutCalled).toBe(true);
      });
    });
  });
});
