import { screen, within, waitFor } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { testApp } from "~test/builder";
import type { Survivor } from "~/lib/survivor";

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

describe("PopulationTab", () => {
  describe("viewing survivors", () => {
    it("shows an empty table when no survivors exist", async () => {
      testApp().renderAt("/settlements/1/population");

      expect(await screen.findByText("No Survivors")).toBeInTheDocument();
    });

    it("shows survivors in the table", async () => {
      testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      expect(await screen.findByText("Test Survivor")).toBeInTheDocument();
    });

    it("displays a row for each survivor", async () => {
      testApp()
        .withSurvivors("1", mockSurvivors)
        .renderAt("/settlements/1/population");

      await screen.findByText("Alice");

      const table = screen.getByRole("table");
      const rows = within(table).getAllByRole("row");
      expect(rows).toHaveLength(4);
    });
  });

  describe("sorting", () => {
    it("sorts by column ascending on first click", async () => {
      const user = testApp()
        .withSurvivors("1", mockSurvivors)
        .renderAt("/settlements/1/population");

      await screen.findByText("Alice");

      await user.click(screen.getByRole("columnheader", { name: /name/i }));

      const table = screen.getByRole("table");
      const rows = within(table).getAllByRole("row");
      const names = rows.slice(1).map(row => within(row).getAllByRole("cell")[0].textContent);
      expect(names).toEqual(["Alice", "Bob", "Charlie"]);
    });

    it("sorts by column descending on second click", async () => {
      const user = testApp()
        .withSurvivors("1", mockSurvivors)
        .renderAt("/settlements/1/population");

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
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      await screen.findByText("Test Survivor");

      expect(screen.queryByRole("columnheader", { name: /birth/i })).not.toBeInTheDocument();

      await user.click(screen.getByTitle("Configure columns"));
      await user.click(screen.getByLabelText("Birth"));

      expect(screen.getByRole("columnheader", { name: /birth/i })).toBeInTheDocument();
    });

    it("can remove a column from the table", async () => {
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      await screen.findByText("Test Survivor");

      expect(screen.getByRole("columnheader", { name: /gender/i })).toBeInTheDocument();

      await user.click(screen.getByTitle("Configure columns"));
      await user.click(screen.getByLabelText("Gender"));

      expect(screen.queryByRole("columnheader", { name: /gender/i })).not.toBeInTheDocument();
    });

    it("cannot remove the name column", async () => {
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      await screen.findByText("Test Survivor");

      await user.click(screen.getByTitle("Configure columns"));

      expect(screen.queryByLabelText("Name")).not.toBeInTheDocument();
    });
  });

  describe("creating a survivor", () => {
    it("disables save when name is invalid", async () => {
      const user = testApp().renderAt("/settlements/1/population");

      await screen.findByText("No Survivors");

      await user.click(screen.getByTitle("Create Survivor"));

      const dialog = screen.getByRole("dialog");
      const nameInput = within(dialog).getByPlaceholderText(/survivor name/i);

      await user.clear(nameInput);

      expect(within(dialog).getByRole("button", { name: /save/i })).toBeDisabled();
    });

    it("updates table after creation", async () => {
      const user = testApp().renderAt("/settlements/1/population");

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
      const user = testApp()
        .withSurvivors("1", [{ ...mockSurvivor, survival: 1 }])
        .renderAt("/settlements/1/population");

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
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");
      expect(dialog).toHaveAttribute("open");

      await user.click(within(dialog).getByRole("button", { name: /cancel/i }));

      expect(dialog).not.toHaveAttribute("open");
    });

    it("does not allow editing the name", async () => {
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

      const survivorRow = await screen.findByText("Test Survivor");

      await user.pointer({ keys: "[MouseRight]", target: survivorRow });
      await user.click(screen.getByRole("button", { name: /edit/i }));

      const dialog = screen.getByRole("dialog");

      expect(within(dialog).queryByPlaceholderText(/survivor name/i)).not.toBeInTheDocument();
      expect(within(dialog).getByText("Test Survivor")).toBeInTheDocument();
    });

    it("does not allow editing the gender", async () => {
      const user = testApp()
        .withSurvivors("1", [mockSurvivor])
        .renderAt("/settlements/1/population");

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
