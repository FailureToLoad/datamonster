import { createMemoryRouter, RouterProvider } from "react-router";
import { GlossaryContext, type GlossaryData } from "~/hooks/glossary";
import { createTestRoutes } from "./routes";

const defaultGlossary: GlossaryData = {
  disorders: [],
  fightingArts: [],
};

export type TestAppOptions = {
  initialPath?: string;
  glossary?: GlossaryData;
};

export function TestApp({ 
  initialPath = "/", 
  glossary = defaultGlossary,
}: TestAppOptions) {
  const router = createMemoryRouter(createTestRoutes(), {
    initialEntries: [initialPath],
  });

  return (
    <GlossaryContext.Provider value={glossary}>
      <RouterProvider router={router} />
    </GlossaryContext.Provider>
  );
}
