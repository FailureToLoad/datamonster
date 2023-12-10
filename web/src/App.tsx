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
    path: "/",
    loader: rootLoader,
    Component: Outlet,
    children: [
      {
        index: true,
        Component: Settlement,
        loader: SettlementLoader,
      },
      {
        path: "/select",
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
