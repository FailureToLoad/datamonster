import { useLoaderData, useRevalidator, Link } from "react-router";
import type { SettlementId } from "~/lib/types/settlement";
import {PlayIcon} from "@phosphor-icons/react";
import { CreateSettlementDialog } from "./CreationDialog";

function SettlementCard({ settlement }: { settlement: SettlementId }) {
  return (
    <div className="card bg-base-100 w-96 overflow-hidden shadow-sm">
      <div className="card-body flex-row items-stretch p-0">
        <div className="flex flex-1 basis-2/3 items-center p-4">
          <h2 className="card-title">{settlement.name}</h2>
        </div>

        <Link
          to={"/settlements/" + settlement.id}
          className="btn btn-ghost h-auto basis-1/3 rounded-none"
        >
          <PlayIcon size={24} />
        </Link>
      </div>
    </div>
  );
}

export default function SettlementsPage() {
  const settlements = useLoaderData() as SettlementId[];
  const revalidator = useRevalidator();

  return (
    <ul className="flex w-full flex-col items-center space-y-4">
      {settlements.map((settlement) => (
        <li key={settlement.id}>
          <SettlementCard settlement={settlement} />
        </li>
      ))}
      <li key="create">
        <CreateSettlementDialog refresh={() => revalidator.revalidate()} />
      </li>
    </ul>
  );
}
