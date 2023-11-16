import { Survivor, columns } from "./columns";
import { DataTable } from "./data-table";
export default function Population() {
  const data: Survivor[] = [
    {
      id: "728ed52f",
      name: "Zach",
      born: 1,
      gender: "M",
      status: "alive",
    },
    {
      id: "728ed52l",
      name: "Lucy",
      born: 1,
      gender: "F",
      status: "alive",
    },
  ];
  return (
    <div id="population">
      <DataTable columns={columns} data={data} />
    </div>
  );
}
