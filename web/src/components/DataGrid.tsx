import { AgGridReact, type AgGridReactProps } from "ag-grid-react";
import { AllCommunityModule, ModuleRegistry } from "ag-grid-community";

ModuleRegistry.registerModules([AllCommunityModule]);

export function DataGrid(props: AgGridReactProps) {
  return (
    <div style={{ width: "100%" }}>
      <AgGridReact
        domLayout="autoHeight"
        overlayNoRowsTemplate="No data"
        tooltipShowDelay={300}
        {...props}
      />
    </div>
  );
}
