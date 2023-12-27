import "./App.css";
import { Authenticator, AuthLoader } from "./auth/authenticator.ts";
import {
  Outlet,
  RouterProvider,
  createBrowserRouter,
  redirect,
} from "react-router-dom";
import Spinner from "./components/spinner.tsx";
import Login, { LoginAction, LoginLoader } from "./routes/login";
import Selector, { SettlementListLoader } from "./routes/settlementSelector";
import { Settlement } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";
import settlementApi from "@/api/settlement.ts";
import survivorApi from "@/api/survivor.ts";

const router = createBrowserRouter([
  {
    id: "root",
    loader: AuthLoader,
    Component: Outlet,
    children: [
      {
        path: ":settlementId",
        id: "home",
        Component: Settlement,
        loader: async ({ params }) => {
          let id = params?.settlementId as string;
          return await settlementApi.getSettlement(id);
        },
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
        Component: Selector,
      },
    ],
  },
  {
    path: "/logout",
    async action() {
      await Authenticator.signout();
      return redirect("/");
    },
  },
  {
    path: "/login",
    action: LoginAction,
    loader: LoginLoader,
    Component: Login,
  },
]);

export default function App() {
  return <RouterProvider router={router} fallbackElement={<Spinner />} />;
}
