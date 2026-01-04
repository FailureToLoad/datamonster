import { redirect } from "react-router";
import type { RouteObject } from "react-router";
import { SelectSettlement, loadSettlements } from "~/pages/selection/index";
import { SettlementPage, TimelineTab, PopulationTab, StorageTab, loadSurvivors } from "~/pages/settlement/index";
import Layout from "~/components/DefaultLayout";

export function createTestRoutes(): RouteObject[] {
  return [
    { path: "/", element: <div>Home</div> },
    {
      element: <Layout />,
      children: [
        {
          path: "/settlements",
          element: <SelectSettlement />,
          loader: loadSettlements,
        },
        {
          path: "/settlements/:settlementId",
          element: <SettlementPage />,
          children: [
            {
              index: true,
              loader: ({ params }) =>
                redirect(`/settlements/${params.settlementId}/timeline`),
            },
            { path: "timeline", element: <TimelineTab /> },
            { 
              path: "population", 
              element: <PopulationTab />,
              loader: loadSurvivors,
            },
            { path: "storage", element: <StorageTab /> },
          ],
        },
      ],
    },
  ];
}
