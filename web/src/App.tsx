import "./App.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Selector, { SettlementListLoader } from "./routes/settlementSelector";
import { Settlement } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";
import settlementApi from "@/api/settlement.ts";
import survivorApi from "@/api/survivor.ts";

import SuperTokens, { SuperTokensWrapper } from "supertokens-auth-react";
import EmailPassword from "supertokens-auth-react/recipe/emailpassword";
import Session from "supertokens-auth-react/recipe/session";
import { getSuperTokensRoutesForReactRouterDom } from "supertokens-auth-react/ui";
import { EmailPasswordPreBuiltUI } from "supertokens-auth-react/recipe/emailpassword/prebuiltui";
import * as reactRouterDom from "react-router-dom";
import Spinner from "./components/spinner.tsx";

SuperTokens.init({
  appInfo: {
    appName: "Data Monster",
    apiDomain: "http://dev.local:8080",
    websiteDomain: "http://dev.local:8090",
    apiBasePath: "/auth",
    websiteBasePath: "/auth",
  },
  recipeList: [EmailPassword.init(), Session.init()],
});

const router = createBrowserRouter([
  ...getSuperTokensRoutesForReactRouterDom(reactRouterDom, [
    EmailPasswordPreBuiltUI,
  ]).map((r) => r.props),
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
]);

export default function App() {
  return (
    <SuperTokensWrapper>
      <RouterProvider router={router} fallbackElement={<Spinner />} />
    </SuperTokensWrapper>
  );
}
