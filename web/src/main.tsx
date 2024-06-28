import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import Auth0ProviderWithNavigate from "@/components/authProvider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Auth0ProviderWithNavigate>
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </Auth0ProviderWithNavigate>
  </React.StrictMode>,
);
