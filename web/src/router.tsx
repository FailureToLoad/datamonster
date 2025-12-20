import { createBrowserRouter, redirect } from "react-router";
import App from "./App";
import Home from "./pages/Home";
import Settlements from "./pages/settlements/Page";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { checkAuth } from "./lib/auth";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorBoundary />,
    children: [
      {
        index: true,
        element: <Home />,
        loader: async () => {
          if (await checkAuth()) return redirect("/settlements");
          return null;
        },
      },
      {
        path: "settlements",
        element: <Settlements />,
        loader: async () => {
          if (!(await checkAuth())) return redirect("/");
          const res = await fetch("/api/settlements", { credentials: "include" });
          if (res.status === 401) return redirect("/");
          if (!res.ok) throw new Response("Failed to load", { status: res.status });
          return res.json();
        },
      },
    ],
  },
]);
