import { createMemoryRouter, RouterProvider } from "react-router";
import type { ReactNode } from "react";
import type { RouteObject } from "react-router";

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

interface TestRouterWithRoutesProps {
  routes: RouteObject[];
  initialEntries?: string[];
}

export function TestRouterWithRoutes({
  routes,
  initialEntries = ["/"],
}: TestRouterWithRoutesProps) {
  const router = createMemoryRouter(routes, { initialEntries });
  return <RouterProvider router={router} />;
}
