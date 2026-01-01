import { useLoaderData, useRevalidator, Link } from "react-router";
import type { SettlementId } from "~/lib/settlement";
import {PlayIcon} from "@phosphor-icons/react";
import { CreateSettlementDialog } from "./dialog";
import styles from "./page.module.css";

function SettlementCard({ settlement }: { settlement: SettlementId }) {
  return (
    <div className={`card bg-base-100 shadow-sm ${styles.settlementCard}`}>
      <div className={`card-body ${styles.cardBody}`}>
        <div className={styles.cardInfo}>
          <h2 className="card-title">{settlement.name}</h2>
        </div>

        <Link
          to={"/settlements/" + settlement.id}
          className={`btn btn-ghost ${styles.playButton}`}
        >
          <PlayIcon size={24} />
        </Link>
      </div>
    </div>
  );
}

export function SelectSettlement() {
  const settlements = useLoaderData() as SettlementId[];
  const revalidator = useRevalidator();

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <CreateSettlementDialog refresh={() => revalidator.revalidate()} />
        <ul className={styles.list}>
          {settlements.map((settlement) => (
            <li key={settlement.id}>
              <SettlementCard settlement={settlement} />
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
