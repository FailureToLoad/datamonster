import { render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TestApp, type TestAppOptions } from "./TestApp";
import { getMockApiState } from "./setup";
import type { SettlementId } from "~/lib/settlement";
import type { Survivor } from "~/lib/survivor";
import type { GlossaryData } from "~/hooks/glossary";
import type { Disorder, FightingArt } from "~/lib/glossary";

class TestAppBuilder {
  private glossary: GlossaryData = { disorders: [], fightingArts: [] };

  withSettlements(settlements: SettlementId[]) {
    getMockApiState().settlements = settlements;
    return this;
  }

  withSurvivors(settlementId: string, survivors: Survivor[]) {
    getMockApiState().survivors.set(settlementId, survivors);
    return this;
  }

  withGlossary(glossary: Partial<GlossaryData>) {
    this.glossary = { ...this.glossary, ...glossary };
    return this;
  }

  withDisorders(disorders: Disorder[]) {
    this.glossary.disorders = disorders;
    return this;
  }

  withFightingArts(fightingArts: FightingArt[]) {
    this.glossary.fightingArts = fightingArts;
    return this;
  }

  unauthenticated() {
    getMockApiState().authenticated = false;
    return this;
  }

  renderAt(path: string) {
    const options: TestAppOptions = {
      initialPath: path,
      glossary: this.glossary,
    };

    render(<TestApp {...options} />);
    return userEvent.setup();
  }
}

export function testApp() {
  return new TestAppBuilder();
}
