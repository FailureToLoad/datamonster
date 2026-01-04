import { createBrowserRouter, redirect, type MiddlewareFunction } from "react-router";
import App from "~/App";
import Home from "~/pages/Home";
import {SelectSettlement, loadSettlements } from "~/pages/selection/index";
import { 
  SettlementPage, 
  TimelineTab, 
  PopulationTab, 
  StorageTab, 
  loadSurvivors  
} from '~/pages/settlement/index.ts';
import Layout from '~/components/DefaultLayout';
import { ErrorBoundary } from "~/components/ErrorBoundary";

async function checkAuth(): Promise<boolean> {
  const response = await fetch("/api/me", { credentials: "include" });
  return response.ok;
}

const authMiddleware: MiddlewareFunction = async (_args, next) => {
  if (!(await checkAuth())) {
    throw redirect("/");
  }
  return next();
};

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
        element: <Layout />,
        middleware: [authMiddleware],
        children: [
          {
            path: '/settlements',
            element: <SelectSettlement />,
            loader: loadSettlements,
          },
          {
            path: '/settlements/:settlementId',
            element: <SettlementPage />,
            children: [
              {
                path: 'population',
                element: <PopulationTab />,
                loader: loadSurvivors,
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
