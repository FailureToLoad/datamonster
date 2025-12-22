import { createBrowserRouter, redirect } from "react-router";
import App from "./App";
import Home from "./pages/Home";
import SettlementsPage from "./pages/settlements/Page";
import SettlementPage from './pages/settlement/Page.tsx';
import PopulationTab from './pages/settlement/population/index.tsx';
import StorageTab from './pages/settlement/SettlementStorage.tsx';
import ProtectedLayout from './components/ProtectedLayout.tsx';
import TimelineTab from './pages/settlement/Timeline.tsx';
import { ErrorBoundary } from "./components/ErrorBoundary";

async function checkAuth(): Promise<boolean> {
  const response = await fetch("/api/me", { credentials: "include" });
  return response.ok;
}

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
        element: <ProtectedLayout />,
        loader: async () => {
          if (!(await checkAuth())) return redirect("/");
          return null;
        },
        children: [
          {
            path: '/settlements',
            element: <SettlementsPage />,
            loader: async () => {
              const res = await fetch("/api/settlements", { credentials: "include" });
              if (res.status === 401) 
                return redirect("/");
              if (!res.ok) 
                throw new Response("Failed to load", { status: res.status });
              return res.json();
          },
          },
          {
            path: '/settlements/:settlementId',
            element: <SettlementPage />,
            children: [
              {
                path: 'population',
                element: <PopulationTab />,
                loader: async ({params}) => {
                  const res = await fetch(`/api/settlements/${params.settlementId}/survivors`, {
                    credentials: 'include',
                  });
                  if (res.status === 401) return redirect('/');
                  if (!res.ok) throw new Response('Failed to load survivors', {status: res.status});
                  return res.json();
                },
              },
              {path: 'storage', element: <StorageTab />},
              {path: 'timeline', element: <TimelineTab />},
            ],
          },
        ],
      },
    ],
  },
]);
