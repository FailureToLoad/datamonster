import {useState} from 'react';
import {SurvivorTable} from './survivorTable';
import EditSurvivorDialog from './editSurvivorDialog';
import {type Survivor} from '~/types/survivor';
import {useLoaderData, useParams, useRevalidator} from 'react-router';


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
    <div id="population" className="flex flex-col w-fill py-4">
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

