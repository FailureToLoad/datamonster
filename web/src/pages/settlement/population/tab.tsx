import { createSignal, Show } from 'solid-js';
import { createAsync, useParams, query, revalidate } from '@solidjs/router';
import { SurvivorTable } from './survivorTable';
import SurvivorDialog from './dialog/index';
import { type Survivor, SurvivorTemplate } from '~/lib/survivor';
import { Get } from '~/lib/request';
import styles from './tab.module.css';

export const getSurvivors = query(async (settlementId: string) => {
  const response = await Get(`/api/settlements/${settlementId}/survivors`);
  if (!response.ok) {
    return [];
  }
  return (await response.json()) as Survivor[];
}, 'survivors');

export default function PopulationTab() {
  const params = useParams<{ settlementId: string }>();
  const survivors = createAsync(() => getSurvivors(params.settlementId));
  const [dialogState, setDialogState] = createSignal<Survivor | null>(null);

  const handleSuccess = () => {
    setDialogState(null);
    revalidate('survivors');
  };

  const handleClose = () => {
    setDialogState(null);
  };

  return (
    <div id="population" class={styles.tab}>
      <SurvivorTable
        data={survivors() ?? []}
        onEditSurvivor={(survivor) => setDialogState(survivor)}
        onCreateSurvivor={() => setDialogState(SurvivorTemplate(params.settlementId))}
      />
      <Show when={dialogState()}>
        {(survivor) => (
          <SurvivorDialog
            data={survivor()}
            settlementId={params.settlementId}
            onClose={handleClose}
            onSuccess={handleSuccess}
          />
        )}
      </Show>
    </div>
  );
}
