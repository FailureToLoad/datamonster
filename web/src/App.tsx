import "./App.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Selector, { SettlementListLoader } from "./routes/settlementSelector";
import { Settlement, SettlementLoader } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";
import survivorApi from "@/api/survivor.ts";
import Spinner from "./components/spinner.tsx";
import AuthGuard from "./components/authGuard";

const router = createBrowserRouter([
  {
    path: ":settlementId",
    id: "home",
    element: <AuthGuard component={Settlement} />,
    loader: SettlementLoader,
    children: [
      {
        path: "timeline",
        Component: Timeline,
      },
      {
        path: "population",
        Component: Population,
        loader: async ({ params }) => {
          let id = params?.settlementId as string;
          return await survivorApi.getSurvivorsForSettlement(id);
        },
      },
      {
        path: "storage",
        Component: SettlementStorage,
      },
    ],
  },
  {
    path: "/",
    loader: SettlementListLoader,
    element: <AuthGuard component={Selector} />,
  },
]);

export default function App() {
  return <RouterProvider router={router} fallbackElement={<Spinner />} />;
}
