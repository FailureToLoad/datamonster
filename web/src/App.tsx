import "./App.css";
import { Authenticator, AuthLoader } from "./auth/authenticator.ts";
import {
  Outlet,
  RouterProvider,
  createBrowserRouter,
  redirect,
} from "react-router-dom";
import Spinner from "./components/spinner.tsx";
import Login, { LoginAction, LoginLoader } from "./routes/login/login.tsx";
import { Selector, SettlementListLoader } from "./routes/settlementSelector";

import { Settlement, SettlementLoader } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";

const router = createBrowserRouter([
  {
    id: "root",
    loader: AuthLoader,
    Component: Outlet,
    children: [
      {
        path: "/",
        Component: Settlement,
        loader: SettlementLoader,
        children: [
          {
            path: "timeline",
            Component: Timeline,
          },
          {
            path: "population",
            Component: Population,
          },
          {
            path: "storage",
            Component: SettlementStorage,
          },
        ],
      },
      {
        path: "select",
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
