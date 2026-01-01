import {useState} from 'react';
import {SurvivorTable} from './survivorTable';
import EditSurvivorDialog from './editSurvivorDialog';
import {type Survivor} from '~/lib/survivor';
import {useLoaderData, useParams, useRevalidator} from 'react-router';
import styles from './tab.module.css';


export function PopulationTab() {
  const {settlementId} = useParams();
  const survivors = (useLoaderData() as Survivor[]) ?? [];
  const revalidator = useRevalidator();
  const [editingSurvivor, setEditingSurvivor] = useState<Survivor | null>(null);

  if (!settlementId) {
    throw Error('settlement id is required');
  }

  const handleEditSuccess = () => {
    setEditingSurvivor(null);
    revalidator.revalidate();
  };

  return (
    <div id="population" className={styles.tab}>
      <SurvivorTable
        data={survivors}
        settlementId={settlementId}
        onEditSurvivor={setEditingSurvivor}
        onSurvivorCreated={() => revalidator.revalidate()}
      />
      <EditSurvivorDialog
        survivor={editingSurvivor}
        onClose={() => setEditingSurvivor(null)}
        onSuccess={handleEditSuccess}
      />
    </div>
  );
}

