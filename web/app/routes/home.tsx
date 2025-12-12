import { redirect } from "react-router";
import type { Route } from "./+types/home";
import { checkAuth, AUTH_LOGIN_URL } from "~/lib/auth.server";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Datamonster" },
    { name: "description", content: "Sign in to Datamonster" },
  ];
}

export async function loader({ request }: Route.LoaderArgs) {
  const isAuthenticated = await checkAuth(request);
  if (isAuthenticated) {
    return redirect("/settlements");
  }
  return { loginUrl: AUTH_LOGIN_URL };
}

export default function Home({ loaderData }: Route.ComponentProps) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background">
      <h1 className="mb-4 text-5xl font-extrabold leading-none tracking-tight">
        Datamonster
      </h1>
      <a href={loaderData.loginUrl}>Sign In</a>
    </div>
  );
}
