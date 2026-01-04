import { createMemoryRouter, RouterProvider } from "react-router";
import type { ReactNode } from "react";

interface TestRouterProps {
  children: ReactNode;
  loaderData?: unknown;
  initialEntries?: string[];
}

export function TestRouter({
  children,
  loaderData,
  initialEntries = ["/"],
}: TestRouterProps) {
  const router = createMemoryRouter(
    [
      {
        path: "*",
        element: children,
        loader: loaderData !== undefined ? () => loaderData : undefined,
      },
    ],
    { initialEntries }
  );

  return <RouterProvider router={router} />;
}
