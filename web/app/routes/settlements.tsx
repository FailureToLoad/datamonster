import type { Route } from "./+types/settlements";
import { requireAuth } from "~/lib/auth.server";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Settlements - Datamonster" },
    { name: "description", content: "Manage your settlements" },
  ];
}

export async function loader({ request }: Route.LoaderArgs) {
  await requireAuth(request);
  return {};
}

export default function Settlements() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background">
      <button className="btn btn-outline" disabled>
        Add Settlement
      </button>
    </div>
  );
}
