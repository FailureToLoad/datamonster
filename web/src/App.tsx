import "./App.css";
import type { LoaderFunctionArgs } from "react-router-dom";
import { Authenticator } from "./auth/authenticator.tsx";
import {
  Outlet,
  RouterProvider,
  createBrowserRouter,
  redirect,
} from "react-router-dom";
import Spinner from "./components/spinner.tsx";
import Login, { LoginAction } from "./routes/login/login.tsx";
import SettlementSelector, {
  settlementListLoader,
} from "./routes/settlementSelector/index.tsx";

import Settlement, { SettlementLoader } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";

const rootLoader = async ({ request }: LoaderFunctionArgs) => {
  const user = await Authenticator.authorize();
  if (!Authenticator.isAuthenticated) {
    let params = new URLSearchParams();
    params.set("from", new URL(request.url).pathname);
    return redirect("/login?" + params.toString());
  }
  return user;
};

const router = createBrowserRouter([
  {
    id: "root",
    loader: rootLoader,
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
        loader: settlementListLoader,
        Component: SettlementSelector,
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
    loader: loginLoader,
    Component: Login,
  },
]);

async function loginLoader() {
  return Authenticator.isAuthenticated;
}

export default function App() {
  return <RouterProvider router={router} fallbackElement={<Spinner />} />;
}
