import { use, type ReactNode } from "react";
import { GlossaryContext, type GlossaryData } from "~/hooks/glossary";
import { Get } from "~/lib/request";

async function fetchGlossary(): Promise<GlossaryData> {
  const res = await Get("/api/glossary");
  if (!res.ok) throw new Error("Failed to load glossary");
  return res.json() as Promise<GlossaryData>;
}

const glossaryPromise = fetchGlossary();

export function GlossaryProvider({ children }: { children: ReactNode }) {
  const glossary = use(glossaryPromise);
  return (
    <GlossaryContext.Provider value={glossary}>
      {children}
    </GlossaryContext.Provider>
  );
}
