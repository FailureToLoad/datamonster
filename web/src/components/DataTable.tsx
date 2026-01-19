import { createMemo, createSignal, For, Show } from 'solid-js';
import {
  createColumnHelper,
  createSolidTable,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  type DeepKeys,
  type SortingState,
  type VisibilityState,
} from '@tanstack/solid-table';
import Settings from 'lucide-solid/icons/settings';
import Plus from 'lucide-solid/icons/plus';
import styles from './DataTable.module.css';

export type ColumnConfig<T> = {
  field: DeepKeys<T>;
  headerName: string;
  headerTooltip?: string;
  default?: boolean;
  sortable?: boolean;
  locked: boolean;
};

export type ContextAction<T> = {
  label: string;
  onClick: (row: T) => void;
};

export type AddRowAction = {
  label: string;
  onClick: () => void;
};

type ContextMenuState<T> = {
  visible: boolean;
  x: number;
  y: number;
  row: T | null;
};

interface DataTableProps<T extends { id: string }> {
  columns: ColumnConfig<T>[];
  rows: T[];
  contextActions?: ContextAction<T>[] | null;
  addRow?: AddRowAction | null;
  emptyMessage?: string;
}

export function DataTable<T extends { id: string }>(props: DataTableProps<T>) {
  const [sorting, setSorting] = createSignal<SortingState>([]);
  const [columnMenuOpen, setColumnMenuOpen] = createSignal(false);
  const contextActions = createMemo(() => props.contextActions ?? []);
  const columnHelper = createColumnHelper<T>();
  const [contextMenu, setContextMenu] = createSignal<ContextMenuState<T>>({
    visible: false,
    x: 0,
    y: 0,
    row: null,
  });

  const initialVisibility = () => {
    const vis: VisibilityState = {};
    for (const col of props.columns) {
      if (!col.locked && !col.default) {
        vis[String(col.field)] = false;
      }
    }
    return vis;
  };

  const [columnVisibility, setColumnVisibility] = createSignal<VisibilityState>(initialVisibility());

  const tanstackColumns = createMemo(() =>
    props.columns.map((col) =>
      columnHelper.accessor(col.field, {
        id: String(col.field),
        header: () => col.headerName,
        cell: (info) => info.getValue(),
        enableSorting: col.sortable ?? false,
        meta: {
          headerTooltip: col.headerTooltip,
          locked: col.locked,
        },
      }),
    ),
  );

  const table = createSolidTable({
    get data() {
      return props.rows;
    },
    get columns() {
      return tanstackColumns();
    },
    getRowId: (row) => row.id,
    state: {
      get sorting() {
        return sorting();
      },
      get columnVisibility() {
        return columnVisibility();
      },
    },
    onSortingChange: setSorting,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  function handleContextMenu(e: MouseEvent, row: T) {
    if (!contextActions().length) return;
    e.preventDefault();
    setContextMenu({
      visible: true,
      x: e.clientX,
      y: e.clientY,
      row,
    });
  }

  function closeContextMenu() {
    setContextMenu((prev) => ({ ...prev, visible: false }));
  }

  const toggleableColumns = () => props.columns.filter((c) => !c.locked);

  return (
    <div
      class={styles.container}
      onContextMenu={(e) => contextActions().length && e.preventDefault()}
    >
      <div class={styles.header}>
        <Show when={props.addRow}>
          {(addRow) => (
            <button
              class={styles.btnGhost}
              aria-label={addRow().label}
              title={addRow().label}
              onClick={addRow().onClick}
            >
              <Plus size={18} strokeWidth={3} />
            </button>
          )}
        </Show>
        <div class={styles.columnMenuWrapper}>
          <button
            class={styles.btnGhost}
            onClick={() => setColumnMenuOpen(!columnMenuOpen())}
            title="Configure columns"
          >
            <Settings size={18} strokeWidth={3} />
          </button>
          <Show when={columnMenuOpen()}>
            <div class={styles.overlay} onClick={() => setColumnMenuOpen(false)} />
            <div class={styles.columnMenu}>
              <p class={styles.columnMenuTitle}>Columns</p>
              <For each={toggleableColumns()}>
                {(column) => (
                  <label class={styles.columnOption}>
                    <input
                      type="checkbox"
                      class="checkbox checkbox-sm"
                      checked={columnVisibility()[String(column.field)] !== false}
                      onChange={() => {
                        const field = String(column.field);
                        setColumnVisibility((prev) => ({
                          ...prev,
                          [field]: prev[field] === false,
                        }));
                      }}
                    />
                    <span class={styles.columnOptionLabel}>{column.headerName}</span>
                  </label>
                )}
              </For>
            </div>
          </Show>
        </div>
      </div>
      <table class={styles.dataTable}>
        <thead>
          <For each={table.getHeaderGroups()}>
            {(headerGroup) => (
              <tr>
                <For each={headerGroup.headers}>
                  {(header) => (
                    <th
                      title={(header.column.columnDef.meta as { headerTooltip?: string })?.headerTooltip}
                      class={header.column.getCanSort() ? styles.sortableHeader : undefined}
                      onClick={header.column.getToggleSortingHandler()}
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(header.column.columnDef.header, header.getContext())}
                      <Show when={header.column.getCanSort()}>
                        <span
                          style={{ opacity: header.column.getIsSorted() ? 1 : 0, 'margin-left': '0.25rem' }}
                        >
                          {header.column.getIsSorted() === 'desc' ? '▼' : '▲'}
                        </span>
                      </Show>
                    </th>
                  )}
                </For>
              </tr>
            )}
          </For>
        </thead>
        <tbody>
          <Show when={table.getRowModel().rows.length === 0 && props.emptyMessage}>
            <tr>
              <td colSpan={table.getVisibleLeafColumns().length} class={styles.emptyMessage}>
                {props.emptyMessage}
              </td>
            </tr>
          </Show>
          <For each={table.getRowModel().rows}>
            {(row) => (
              <tr onContextMenu={(e) => handleContextMenu(e, row.original)}>
                <For each={row.getVisibleCells()}>
                  {(cell) => (
                    <td>{flexRender(cell.column.columnDef.cell, cell.getContext())}</td>
                  )}
                </For>
              </tr>
            )}
          </For>
        </tbody>
      </table>
      <Show when={contextMenu().visible && contextActions().length}>
        <div
          class={styles.overlay}
          onClick={closeContextMenu}
          onContextMenu={(e) => {
            e.preventDefault();
            closeContextMenu();
          }}
        />
        <div class={styles.contextMenu} style={{ left: `${contextMenu().x}px`, top: `${contextMenu().y}px` }}>
          <For each={contextActions()}>
            {(action) => (
              <button
                class={styles.contextMenuItem}
                onClick={() => {
                  const row = contextMenu().row;
                  if (row) action.onClick(row);
                  closeContextMenu();
                }}
              >
                {action.label}
              </button>
            )}
          </For>
        </div>
      </Show>
    </div>
  );
}
