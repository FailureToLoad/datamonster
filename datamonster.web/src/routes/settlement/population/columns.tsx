import { ColumnDef } from "@tanstack/react-table";
import { ArrowDown, ArrowUp } from "lucide-react";
import { Button } from "@/components/ui/button";

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Survivor = {
  id: string;
  name: string;
  born: number;
  gender: "M" | "F";
  status: "alive" | "dead" | "retired";
};

export const columns: ColumnDef<Survivor>[] = [
  {
    accessorKey: "name",
    header: ({ column }) => {
      const isAsc = column.getIsSorted() === "asc";
      return (
        <Button variant="ghost" onClick={() => column.toggleSorting(isAsc)}>
          Name
          {isAsc ? (
            <ArrowDown className="ml-2 h-4 w-4" />
          ) : (
            <ArrowUp className="ml-2 h-4 w-4" />
          )}
        </Button>
      );
    },
  },
  {
    accessorKey: "born",
    header: "Born",
  },
  {
    accessorKey: "gender",
    header: "Gender",
  },
  {
    accessorKey: "status",
    header: "Status",
  },
];
