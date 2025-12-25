import { DataGrid } from '~/components/DataGrid';
import { type Survivor } from '~/types/survivor';

export function SurvivorTable({ data }: { data: Survivor[] }) {
    return (
        <DataGrid
            rowData={data}
            defaultColDef={{ cellStyle: { textAlign: "center" }, headerClass: "ag-header-cell-center" }}
            columnDefs={[
                { field: "name", cellStyle: { textAlign: "left" }, headerClass: "" },
                { field: "gender" },
                { field: "born" },
                { field: "huntxp" },
                { field: "survival" },
                { field: "movement" },
                { field: "speed" },
                { field: "strength" },
                { field: "accuracy" },
                { field: "evasion" },
                { field: "luck" },
                { field: "systemicpressure", headerName: "S.P.", headerTooltip: "System Pressure" },
                { field: "torment" },
                { field: "courage" },
                { field: "understanding" },
                { field: "insanity" },
                { field: "lumi" },
            ]}
            autoSizeStrategy={{ type: "fitGridWidth" }}
        />
    );
}