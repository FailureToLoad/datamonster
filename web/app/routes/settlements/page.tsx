import { redirect, useLoaderData, useRevalidator } from "react-router";
import type { SettlementId } from "~/lib/types/settlement";
import Play from "lucide-react/dist/esm/icons/play";
import { Link } from "react-router";
import type { Route } from "./+types/page";
import { checkAuth } from "~/lib/auth.server";
import { CreateSettlementDialog } from "./creationDialog";

const API_URL = process.env.API_URL;

export async function loader({ request }: Route.LoaderArgs) {
  const isAuthenticated = await checkAuth(request);
  if (!isAuthenticated) {
    throw redirect("/");
  }

  const response = await fetch(`${API_URL}/api/settlements`, {
    headers: {
      cookie: request.headers.get("cookie") || "",
    },
  });

  if (!response.ok) {
    throw new Response("Failed to load settlements", { status: response.status });
  }

  return response.json();
}

function SettlementCard({settlement}: {settlement: SettlementId}) {
  return (
    <div className="card bg-base-100 w-96 shadow-sm">
      <div className="card-body">
        <h2 className="card-title">{settlement.name}</h2>
        <div className="flex flex-row justify-between">
          <div>
            <Link to={'/settlements/' + settlement.id}>
              <button className="btn btn-ghost">
                <Play size={24} />
              </button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function SettlementsPage() {
  const settlements = useLoaderData<typeof loader>() as SettlementId[];
  const revalidator = useRevalidator();

  return (
    <main className="flex w-screen h-screen flex-col items-center justify-center overflow-hidden">
      <ul className="w-1/4 space-y-4">
        {settlements.map((settlement) => (
          <li key={settlement.id}>
            <SettlementCard settlement={settlement} />
          </li>
        ))}
        <li key="create">
          <CreateSettlementDialog refresh={() => revalidator.revalidate()} />
        </li>
      </ul>
    </main>
  );
}
