import "./App.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Selector from "@/routes/settlementSelector";
import { Settlement } from "./routes/settlement";
import Timeline from "./routes/settlement/timeline.tsx";
import Population from "./routes/settlement/population/index.tsx";
import SettlementStorage from "./routes/settlement/settlementStorage.tsx";
import Spinner from "./components/spinner.tsx";
import AuthGuard from "./components/authGuard";
import Welcome from "./routes/welcome/index.tsx";
import Layout from "./layout.tsx";
import ErrorPage from "./routes/error.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    id: "root",
    element: <Layout />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: ":settlementId",
        id: "home",
        element: <AuthGuard component={Settlement} />,
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
        path: "/select",
        element: <AuthGuard component={Selector} />,
      },
      {
        path: "/welcome",
        Component: Welcome,
      },
    ],
  },
]);

export default function App() {
  return <RouterProvider router={router} fallbackElement={<Spinner />} />;
}
