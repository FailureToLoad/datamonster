import { A, query, createAsync, revalidate } from "@solidjs/router";
import { For } from "solid-js";
import Play from "lucide-solid/icons/play";
import type { SettlementId } from "~/lib/settlement";
import { CreateSettlementDialog } from "./dialog";
import { Get } from "~/lib/request";
import styles from "./page.module.css";

export const getSettlements = query(async (): Promise<SettlementId[]> => {
  const response = await Get("/api/settlements");
  return response.json();
}, "settlements");

function SettlementCard(props: { settlement: SettlementId }) {
  return (
    <div class={styles.settlementCard}>
      <div class={styles.cardBody}>
        <div class={styles.cardInfo}>
          <h2 class="card-title">{props.settlement.name}</h2>
        </div>

        <A
          href={"/settlements/" + props.settlement.id}
          class={styles.playButton}
        >
          <Play size={24} />
        </A>
      </div>
    </div>
  );
}

export default function SelectSettlement() {
  const settlements = createAsync(() => getSettlements());

  const refresh = () => revalidate(getSettlements.key);

  return (
    <div class={styles.container}>
      <div class={styles.content}>
        <CreateSettlementDialog refresh={refresh} />
        <ul class={styles.list}>
          <For each={settlements() ?? []}>
            {(settlement) => (
              <li>
                <SettlementCard settlement={settlement} />
              </li>
            )}
          </For>
        </ul>
      </div>
    </div>
  );
}
