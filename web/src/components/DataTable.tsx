import { useMemo, useState } from "react";
import { GearIcon, PlusIcon } from "@phosphor-icons/react";
import styles from "./DataTable.module.css";

export type ColumnConfig<T> = {
    field: keyof T;
    headerName: string;
    headerTooltip?: string;
    default?: boolean;
    sortable?: boolean;
    locked: boolean;
};

type SortState<T> = {
    field: keyof T;
    direction: "asc" | "desc";
} | null;

type ContextMenuState<T> = {
    visible: boolean;
    x: number;
    y: number;
    row: T | null;
};

export type ContextAction<T> = {
    label: string;
    onClick: (row: T) => void;
};

export type AddRowAction = {
    label: string;
    onClick: () => void;
};

type DataTableProps<T extends { id: string }> = {
    columns: ColumnConfig<T>[];
    rows: T[];
    contextActions?: ContextAction<T>[] | null;
    addRow?: AddRowAction | null;
};

export function DataTable<T extends { id: string }>({ columns, rows, contextActions, addRow }: DataTableProps<T>) {
    const [sort, setSort] = useState<SortState<T>>(null);
    const [columnMenuOpen, setColumnMenuOpen] = useState(false);
    const [contextMenu, setContextMenu] = useState<ContextMenuState<T>>({
        visible: false,
        x: 0,
        y: 0,
        row: null,
    });
    const [visibleFields, setVisibleFields] = useState<Set<keyof T>>(() =>
        new Set(columns.filter(c => c.default).map(c => c.field))
    );

    const visibleColumns = useMemo(() =>
        columns.filter(c => c.locked || visibleFields.has(c.field)),
        [columns, visibleFields]
    );

    function toggleColumn(field: keyof T) {
        setVisibleFields(prev => {
            const next = new Set(prev);
            if (next.has(field)) {
                next.delete(field);
            } else {
                next.add(field);
            }
            return next;
        });
    }

    const sortedRows = useMemo(() => {
        if (!sort) return rows;
        return [...rows].sort((a, b) => {
            const aVal = a[sort.field];
            const bVal = b[sort.field];
            if (aVal < bVal) return sort.direction === "asc" ? -1 : 1;
            if (aVal > bVal) return sort.direction === "asc" ? 1 : -1;
            return 0;
        });
    }, [rows, sort]);

    function handleSort(field: keyof T) {
        setSort((prev: SortState<T>) => {
            if (prev?.field !== field) return { field, direction: "asc" };
            if (prev.direction === "asc") return { field, direction: "desc" };
            return null;
        });
    }

    function handleContextMenu(e: React.MouseEvent, row: T) {
        if (!contextActions?.length) return;
        e.preventDefault();
        setContextMenu({
            visible: true,
            x: e.clientX,
            y: e.clientY,
            row,
        });
    }

    function closeContextMenu() {
        setContextMenu((prev: ContextMenuState<T>) => ({ ...prev, visible: false }));
    }

    return (
        <div className={styles.container} onContextMenu={(e) => contextActions?.length && e.preventDefault()}>
            <div className={styles.header}>
                {addRow && (
                    <button
                        className={styles.btnGhost}
                        aria-label={addRow.label}
                        title={addRow.label}
                        onClick={addRow.onClick}
                    >
                        <PlusIcon size={18} weight="bold" />
                    </button>
                )}
                <div className={styles.columnMenuWrapper}>
                    <button
                        className={styles.btnGhost}
                        onClick={() => setColumnMenuOpen(!columnMenuOpen)}
                        title="Configure columns"
                    >
                        <GearIcon size={18} weight="bold" />
                    </button>
                    {columnMenuOpen && (
                        <>
                            <div
                                className={styles.overlay}
                                onClick={() => setColumnMenuOpen(false)}
                            />
                            <div className={styles.columnMenu}>
                                <p className={styles.columnMenuTitle}>Columns</p>
                                {columns.filter(c => !c.locked).map((column) => (
                                    <label
                                        key={String(column.field)}
                                        className={styles.columnOption}
                                    >
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm"
                                            checked={visibleFields.has(column.field)}
                                            onChange={() => toggleColumn(column.field)}
                                        />
                                        <span className={styles.columnOptionLabel}>{column.headerName}</span>
                                    </label>
                                ))}
                            </div>
                        </>
                    )}
                </div>
            </div>
            <table className={styles.dataTable}>
                <thead>
                    <tr>
                        {visibleColumns.map((column) => (
                            <th
                                key={String(column.field)}
                                title={column.headerTooltip}
                                onClick={column.sortable ? () => handleSort(column.field) : undefined}
                                style={column.sortable ? { cursor: "pointer" } : undefined}
                            >
                                {column.headerName}
                                {column.sortable && (
                                    <span style={{ opacity: sort?.field === column.field ? 1 : 0 }}>
                                        {" "}{sort?.field === column.field && sort.direction === "desc" ? "▼" : "▲"}
                                    </span>
                                )}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {sortedRows.map((row) => (
                        <tr
                            key={row.id}
                            onContextMenu={(e) => handleContextMenu(e, row)}
                        >
                            {visibleColumns.map((column) => (
                                <td key={String(column.field)}>
                                    {String(row[column.field])}
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
            {contextMenu.visible && contextActions?.length && (
                <>
                    <div
                        className={styles.overlay}
                        onClick={closeContextMenu}
                        onContextMenu={(e) => { e.preventDefault(); closeContextMenu(); }}
                    />
                    <div
                        className={styles.contextMenu}
                        style={{ left: contextMenu.x, top: contextMenu.y }}
                    >
                        {contextActions.map((action) => (
                            <button
                                key={action.label}
                                className={styles.contextMenuItem}
                                onClick={() => {
                                    if (contextMenu.row) action.onClick(contextMenu.row);
                                    closeContextMenu();
                                }}
                            >
                                {action.label}
                            </button>
                        ))}
                    </div>
                </>
            )}
        </div>
    )
}