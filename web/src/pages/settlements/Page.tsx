import { useLoaderData, useRevalidator, Link } from "react-router";
import type { SettlementId } from "~/lib/types/settlement";
import Play from "lucide-react/dist/esm/icons/play";
import { CreateSettlementDialog } from "./CreationDialog";

function SettlementCard({ settlement }: { settlement: SettlementId }) {
  return (
    <div className="card bg-base-100 w-96 shadow-sm">
      <div className="card-body flex-row items-center p-0">
        <div className="flex-1 basis-2/3 p-4">
          <h2 className="card-title">{settlement.name}</h2>
        </div>
        <Link
          to={"/settlements/" + settlement.id}
          className="basis-1/3 py-6 flex items-center justify-center border-l border-base-300 hover:bg-base-200 transition-colors"
        >
          <Play size={24} />
        </Link>
      </div>
    </div>
  );
}

export default function SettlementsPage() {
  const settlements = useLoaderData() as SettlementId[];
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
