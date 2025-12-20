import { Outlet, ScrollRestoration } from "react-router";

export default function App() {
  return (
    <>
      <Outlet />
      <ScrollRestoration />
    </>
  );
}
