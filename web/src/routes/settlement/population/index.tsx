import { columns } from "./columns";
import { DataTable } from "./data-table";
import { Survivor } from "./types";
export default function Population() {
  const data: Survivor[] = [
    {
      id: "728ed52f",
      name: "Zach",
      born: 1,
      gender: "M",
      status: "alive",
      xp: 0,
      survival: 1,
      movement: 5,
      accuracy: 0,
      strength: 0,
      evasion: 0,
      luck: 0,
      speed: 0,
      insanity: 0,
      sp: 0,
      torment: 0,
      lumi: 0,
      courage: 0,
      understanding: 0,
    },
    {
      id: "728ed52l",
      name: "Lucy",
      born: 1,
      gender: "F",
      status: "alive",
      xp: 0,
      survival: 1,
      movement: 5,
      accuracy: 0,
      strength: 1,
      evasion: 0,
      luck: 0,
      speed: 0,
      insanity: 0,
      sp: 0,
      torment: 0,
      lumi: 0,
      courage: 0,
      understanding: 0,
    },
  ];
  return (
    <div id="population">
      <DataTable columns={columns} data={data} />
    </div>
  );
}
