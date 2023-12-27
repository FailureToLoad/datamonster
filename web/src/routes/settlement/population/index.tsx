import { useLoaderData } from "react-router-dom";
import { columns } from "./columns";
import { DataTable } from "./data-table";
import { Survivor } from "./types";

export default function Population() {
  const data = useLoaderData() as Survivor[];
  return (
    <div id="population">
      <DataTable columns={columns} data={data} />
    </div>
  );
}
