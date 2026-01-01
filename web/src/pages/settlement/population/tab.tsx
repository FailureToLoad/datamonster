import {useState} from 'react';
import {SurvivorTable} from './survivorTable';
import SurvivorDialog from './survivorDialog';
import {type Survivor, SurvivorTemplate} from '~/lib/survivor';
import {useLoaderData, useParams, useRevalidator} from 'react-router';
import styles from './tab.module.css';

type DialogState = Survivor | null;

export function PopulationTab() {
  const {settlementId} = useParams();
  const survivors = (useLoaderData() as Survivor[]) ?? [];
  const revalidator = useRevalidator();
  const [dialogState, setDialogState] = useState<DialogState>(null);

  if (!settlementId) {
    throw Error('settlement id is required');
  }

  const handleSuccess = () => {
    setDialogState(null);
    revalidator.revalidate();
  };

  const handleClose = () => {
    setDialogState(null);
  };

  return (
    <div id="population" className={styles.tab}>
      <SurvivorTable
        data={survivors}
        onEditSurvivor={(survivor) => setDialogState(survivor)}
        onCreateSurvivor={() => setDialogState(SurvivorTemplate(settlementId))}
      />
      {dialogState && (
        <SurvivorDialog
          data={dialogState}
          settlementId={settlementId}
          onClose={handleClose}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  );
}
