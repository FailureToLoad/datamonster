import type { Accessor } from "solid-js";
import { createContext, useContext } from "solid-js";
import type { Disorder, FightingArt } from "~/lib/glossary";

export type GlossaryData = {
  disorders: Disorder[];
  fightingArts: FightingArt[];
};

export const GlossaryContext = createContext<Accessor<GlossaryData> | null>(null);

export function useGlossary() {
  const context = useContext(GlossaryContext);
  if (!context) {
    throw new Error("useGlossary must be used within a GlossaryProvider");
  }
  return context;
}
