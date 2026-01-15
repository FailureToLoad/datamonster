import { render, screen, within, waitFor } from "@solidjs/testing-library";
import { MemoryRouter, Route, createMemoryHistory, revalidate } from "@solidjs/router";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { Suspense, createRoot } from "solid-js";
import { server } from "~test/setup";
import type { Survivor } from "~/lib/survivor";
import { GlossaryContext, type GlossaryData } from "~/hooks/glossary";
import PopulationTab, { getSurvivors } from "./tab";

const mockGlossary: GlossaryData = {
  disorders: [
    { id: "d1", name: "Anxiety", source: "Core", effect: "Cannot use fighting arts" },
    { id: "d2", name: "Fear of the Dark", source: "Core", effect: "-1 accuracy in darkness" },
  ],
  fightingArts: [
    { id: "fa1", name: "Sword Mastery", secret: false, source: "Core", text: ["Gain +1 accuracy with swords"] },
    { id: "fa2", name: "Red Fist", secret: true, source: "Core", text: ["Secret fighting art effect"] },
  ],
};

const mockSurvivor: Survivor = {
  id: "1",
  name: "Test Survivor",
  gender: "M",
  birth: 1,
  status: "Alive",
  huntxp: 2,
  movement: 5,
  speed: 0,
  strength: 0,
  accuracy: 0,
  evasion: 0,
  luck: 0,
  systemicPressure: 0,
  torment: 0,
  courage: 3,
  understanding: 2,
  survival: 1,
  insanity: 0,
  lumi: 0,
  settlementId: "1",
  disorders: [],
  fightingArt: null,
  secretFightingArt: null,
};

const mockSurvivors: Survivor[] = [
  { ...mockSurvivor, id: "1", name: "Alice", huntxp: 5 },
  { ...mockSurvivor, id: "2", name: "Bob", huntxp: 2 },
  { ...mockSurvivor, id: "3", name: "Charlie", huntxp: 8 },
];

function renderPopulationTab(survivors: Survivor[] = []) {
  server.use(
    http.get("/api/settlements/:settlementId/survivors", () => {
      return HttpResponse.json(survivors);
    }),
    http.post("/api/settlements/:settlementId/survivors", async ({ request }) => {
      const body = (await request.json()) as Survivor;
      const newSurvivor = { ...body, id: "new-id" };
      survivors.push(newSurvivor);
      return HttpResponse.json(newSurvivor);
    }),
    http.patch("/api/settlements/:settlementId/survivors/:id", async ({ request, params }) => {
      const body = (await request.json()) as { statUpdates?: Record<string, number> };
      const survivor = survivors.find((s) => s.id === params.id);
      if (!survivor) return HttpResponse.json({ error: "Not found" }, { status: 404 });

      if (body.statUpdates) Object.assign(survivor, body.statUpdates);
      return HttpResponse.json(survivor);
    })
  );

  function Wrapper() {
    const history = createMemoryHistory();
    history.set({ value: "/settlements/1/population" });

    return (
      <GlossaryContext.Provider value={() => mockGlossary}>
        <MemoryRouter
          history={history}
          root={(props) => (
            <Suspense fallback={<div>Loading...</div>}>{props.children}</Suspense>
          )}
        >
          <Route
            path="/settlements/:settlementId/population"
            component={PopulationTab}
          />
        </MemoryRouter>
      </GlossaryContext.Provider>
    );
  }

  const renderResult = render(() => <Wrapper />);

  return {
    user: userEvent.setup(),
    ...renderResult
  };
}

describe("PopulationTab", () => {
  beforeEach(() => {
    createRoot((dispose) => {
      revalidate(getSurvivors.key, true);
      dispose();
    });
    
    server.resetHandlers();
  });

  afterEach(() => {
    server.resetHandlers();
  });

  describe("viewing survivors", () => {
    it("shows an empty table when no survivors exist", async () => {
      renderPopulationTab([]);

      expect(await screen.findByText("No Survivors")).toBeInTheDocument();
    });

    it("shows survivors in the table", async () => {
      renderPopulationTab([mockSurvivor]);

      expect(await screen.findByText("Test Survivor")).toBeInTheDocument();
    });

    it("displays a row for each survivor", async () => {
      renderPopulationTab(mockSurvivors);

      await screen.findByText("Alice");

      const table = screen.getByRole("table");
      const rows = within(table).getAllByRole("row");
      expect(rows).toHaveLength(4);
    });
  });

  describe("sorting", () => {
    it("sorts by column ascending on first click", async () => {
      const { user } = renderPopulationTab(mockSurvivors);

      await screen.findByText("Alice");

      await user.click(screen.getByRole("columnheader", { name: /name/i }));

      const table = screen.getByRole("table");
      const rows = within(table).getAllByRole("row");
      const names = rows.slice(1).map(row => within(row).getAllByRole("cell")[0].textContent);
      expect(names).toEqual(["Alice", "Bob", "Charlie"]);
    });

    it("sorts by column descending on second click", async () => {
      const { user } = renderPopulationTab(mockSurvivors);

      await screen.findByText("Alice");

      await user.click(screen.getByRole("columnheader", { name: /name/i }));
      await user.click(screen.getByRole("columnheader", { name: /name/i }));

      const table = screen.getByRole("table");
      const rows = within(table).getAllByRole("row");
      const names = rows.slice(1).map(row => within(row).getAllByRole("cell")[0].textContent);
      expect(names).toEqual(["Charlie", "Bob", "Alice"]);
    });
  });

  describe("configuring columns", () => {
    it("can add a column to the table", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      await screen.findByText("Test Survivor");

      expect(screen.queryByRole("columnheader", { name: /birth/i })).not.toBeInTheDocument();

      await user.click(screen.getByTitle("Configure columns"));
      await user.click(screen.getByLabelText("Birth"));

      expect(screen.getByRole("columnheader", { name: /birth/i })).toBeInTheDocument();
    });

    it("can remove a column from the table", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      await screen.findByText("Test Survivor");

      expect(screen.getByRole("columnheader", { name: /gender/i })).toBeInTheDocument();

      await user.click(screen.getByTitle("Configure columns"));
      await user.click(screen.getByLabelText("Gender"));

      expect(screen.queryByRole("columnheader", { name: /gender/i })).not.toBeInTheDocument();
    });

    it("cannot remove the name column", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      await screen.findByText("Test Survivor");

      await user.click(screen.getByTitle("Configure columns"));

      expect(screen.queryByLabelText("Name")).not.toBeInTheDocument();
    });
  });

  describe("creating a survivor", () => {
    it("disables save when name is invalid", async () => {
      const { user } = renderPopulationTab([]);

      await screen.findByText("No Survivors");

      await user.click(screen.getByTitle("Create Survivor"));

      const dialog = screen.getByRole("dialog");
      const nameInput = within(dialog).getByPlaceholderText(/survivor name/i);

      await user.clear(nameInput);

      expect(within(dialog).getByRole("button", { name: /save/i })).toBeDisabled();
    });

    it("updates table after creation", async () => {
      const { user } = renderPopulationTab([]);

      await screen.findByText("No Survivors");

      await user.click(screen.getByTitle("Create Survivor"));

      const dialog = screen.getByRole("dialog");
      const nameInput = within(dialog).getByPlaceholderText(/survivor name/i);

      await user.clear(nameInput);
      await user.type(nameInput, "New Survivor");
      await user.click(within(dialog).getByRole("button", { name: /save/i }));

      await waitFor(() => {
        expect(dialog).not.toHaveAttribute("open");
      });

      expect(await screen.findByText("New Survivor")).toBeInTheDocument();
    });
  });

  describe("editing a survivor", () => {
    it("opens edit dialog on right-click and updates row after save", async () => {
      const { user } = renderPopulationTab([{ ...mockSurvivor, survival: 1 }]);

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");
      expect(dialog).toHaveAttribute("open");

      const survivalInputs = within(dialog).getAllByRole("spinbutton");
      const survivalInputEl = survivalInputs.find(el => el.id === "survival-edit-input");
      expect(survivalInputEl).toBeDefined();
      await user.clear(survivalInputEl!);
      await user.type(survivalInputEl!, "5");

      await user.click(within(dialog).getByRole("button", { name: /save/i }));

      await waitFor(() => {
        expect(dialog).not.toHaveAttribute("open");
      });

      const table = screen.getByRole("table");
      expect(within(table).getByText("5")).toBeInTheDocument();
    });

    it("closes dialog when clicking cancel", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");
      expect(dialog).toHaveAttribute("open");

      await user.click(within(dialog).getByRole("button", { name: /cancel/i }));

      expect(dialog).not.toHaveAttribute("open");
    });

    it("does not allow editing the name", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");

      expect(within(dialog).queryByPlaceholderText(/survivor name/i)).not.toBeInTheDocument();
      expect(within(dialog).getByText("Test Survivor")).toBeInTheDocument();
    });

    it("does not allow editing the gender", async () => {
      const { user } = renderPopulationTab([mockSurvivor]);

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");

      const genderButtons = within(dialog).queryAllByRole("button").filter(
        btn => btn.querySelector("svg")
      );
      const clickableGenderButton = genderButtons.find(btn =>
        btn.classList.contains("genderBadgeMale") || btn.classList.contains("genderBadgeFemale")
      );
      expect(clickableGenderButton).toBeUndefined();
    });
  });
});
