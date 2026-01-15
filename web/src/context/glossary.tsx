import { createAsync, query } from "@solidjs/router";
import type { ParentProps } from "solid-js";
import { GlossaryContext, type GlossaryData } from "~/hooks/glossary";
import { Get } from "~/lib/request";

export const getGlossary = query(async () => {
  const res = await Get("/api/glossary");
  if (!res.ok) {
    throw new Error("Failed to load glossary");
  }
  return (await res.json()) as GlossaryData;
}, "glossary");

export function GlossaryProvider(props: ParentProps) {
  const glossary = createAsync(() => getGlossary(), {
    initialValue: { disorders: [], fightingArts: [] },
  });

  return (
    <GlossaryContext.Provider value={glossary}>
      {props.children}
    </GlossaryContext.Provider>
  );
}
