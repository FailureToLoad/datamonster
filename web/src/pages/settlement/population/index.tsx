import {SurvivorTable} from './survivorTable';
import NewSurvivorDialog from './survivorDialog';
import {type Survivor} from '~/types/survivor';
import {useLoaderData, useParams, useRevalidator} from 'react-router';

export default function PopulationTab() {
  const {settlementId} = useParams();
  const survivors = useLoaderData() as Survivor[];
  const revalidator = useRevalidator();

  if (!settlementId) {
    throw Error('settlement id is required');
  }

  return (
    <div id="population" className="max-w-fit py-4">
      <div className="flex flex-row-reverse items-center py-4">
        <NewSurvivorDialog
          settlementId={settlementId}
          onSuccess={() => revalidator.revalidate()}
        />
      </div>
      <SurvivorTable data={survivors} />
    </div>
  );
}
