import {type Survivor} from '~/types/survivor';

interface DataTableProps {
  data: Survivor[];
}

export function SurvivorTable({data}: DataTableProps) {
  return (
    <div className="overflow-x-auto h-96">
      <table className="table table-xs table-pin-rows">
        <thead>
          <tr>
            <td>Name</td>
            <td>Gender</td>
            <td>XP</td>
            <td>Survival</td>
            <td>Movement</td>
            <td>Accuracy</td>
            <td>Strength</td>
            <td>Evasion</td>
            <td>Luck</td>
            <td>Speed</td>
            <td>Insanity</td>
            <td>SP</td>
            <td>Torment</td>
            <td>Lumi</td>
            <td>Courage</td>
            <td>Understanding</td>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {data && data.length > 0 ? (
            data.map((survivor) => (
              <tr key={survivor.id} className="hover">
                <th>{survivor.name}</th>
                <td>{survivor.gender}</td>
                <td>{survivor.huntxp}</td>
                <td>{survivor.survival}</td>
                <td>{survivor.movement}</td>
                <td>{survivor.accuracy}</td>
                <td>{survivor.strength}</td>
                <td>{survivor.evasion}</td>
                <td>{survivor.luck}</td>
                <td>{survivor.speed}</td>
                <td>{survivor.insanity}</td>
                <td>{survivor.systemicpressure}</td>
                <td>{survivor.torment}</td>
                <td>{survivor.lumi}</td>
                <td>{survivor.courage}</td>
                <td>{survivor.understanding}</td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={17} className="h-24 text-center">
                No survivors.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
