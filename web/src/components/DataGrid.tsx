import { AgGridReact, type AgGridReactProps } from "ag-grid-react";
import { AllCommunityModule, ModuleRegistry } from "ag-grid-community";
import styles from "./DataGrid.module.css";

ModuleRegistry.registerModules([AllCommunityModule]);

export function DataGrid(props: AgGridReactProps) {
  return (
    <div className={styles.container}>
      <AgGridReact
        domLayout="autoHeight"
        overlayNoRowsTemplate="No data"
        tooltipShowDelay={300}
        {...props}
      />
    </div>
  );
}
