import {useState, useRef, useCallback} from 'react';
import {DataGrid} from '~/components/DataGrid';
import {type Survivor} from '~/lib/survivor';
import type {CellContextMenuEvent, GridApi, ColDef} from 'ag-grid-community';
import {GearIcon, PlusIcon} from '@phosphor-icons/react';
import styles from './survivorTable.module.css';

const COLUMN_CONFIG_KEY = 'survivor-table-columns';

type SurvivorTableProps = {
    data: Survivor[];
    onEditSurvivor?: (survivor: Survivor) => void;
    onCreateSurvivor?: () => void;
};

type ContextMenuState = {
    visible: boolean;
    x: number;
    y: number;
    survivor: Survivor | null;
};

type ColumnConfig = {
    field: keyof Survivor;
    headerName: string;
    headerTooltip?: string;
    locked?: boolean;
    default?: boolean;
};

const columnConfigs: ColumnConfig[] = [
    {field: 'name', headerName: 'Name', locked: true, default: true},
    {field: 'gender', headerName: 'Gender', default: true},
    {field: 'birth', headerName: 'Birth'},
    {field: 'huntxp', headerName: 'XP', headerTooltip: 'Hunt XP', default: true},
    {field: 'survival', headerName: 'Survival', default: true},
    {field: 'movement', headerName: 'Movement'},
    {field: 'speed', headerName: 'Speed'},
    {field: 'strength', headerName: 'Strength'},
    {field: 'accuracy', headerName: 'Accuracy'},
    {field: 'evasion', headerName: 'Evasion'},
    {field: 'luck', headerName: 'Luck'},
    {field: 'systemicPressure', headerName: 'S.P.', headerTooltip: 'Systemic Pressure'},
    {field: 'torment', headerName: 'Torment'},
    {field: 'courage', headerName: 'Courage', default: true},
    {field: 'understanding', headerName: 'Understanding', default: true},
    {field: 'insanity', headerName: 'Insanity', default: true},
    {field: 'lumi', headerName: 'Lumi'},
];

export function SurvivorTable({data, onEditSurvivor, onCreateSurvivor}: SurvivorTableProps) {
    const gridApiRef = useRef<GridApi | null>(null);
    const [columnMenuOpen, setColumnMenuOpen] = useState(false);
    const [visibleColumns, setVisibleColumns] = useState<Set<string>>(() => {
        const stored = localStorage.getItem(COLUMN_CONFIG_KEY);
        if (stored) {
            try {
                return new Set(JSON.parse(stored) as string[]);
            } catch {
                return new Set(columnConfigs.filter(col => col.default).map(col => col.field));
            }
        }
        return new Set(columnConfigs.filter(col => col.default).map(col => col.field));
    });

    const [contextMenu, setContextMenu] = useState<ContextMenuState>({
        visible: false,
        x: 0,
        y: 0,
        survivor: null,
    });

    const closeMenu = () => {
        gridApiRef.current?.deselectAll();
        setContextMenu(prev => ({...prev, visible: false}));
    };

    const handleCellContextMenu = (event: CellContextMenuEvent<Survivor>) => {
        event.event?.preventDefault();
        const mouseEvent = event.event as MouseEvent;
        if (event.data && event.node) {
            event.node.setSelected(true);
            setContextMenu({
                visible: true,
                x: mouseEvent.clientX,
                y: mouseEvent.clientY,
                survivor: event.data,
            });
        }
    };

    const handleEdit = () => {
        if (contextMenu.survivor && onEditSurvivor) {
            onEditSurvivor(contextMenu.survivor);
        }
        closeMenu();
    };

    const toggleColumn = useCallback((field: string) => {
        const config = columnConfigs.find(c => c.field === field);
        if (config?.locked) return;

        setVisibleColumns(prev => {
            const next = new Set(prev);
            if (next.has(field)) {
                next.delete(field);
            } else {
                next.add(field);
            }
            gridApiRef.current?.setColumnsVisible([field], next.has(field));
            localStorage.setItem(COLUMN_CONFIG_KEY, JSON.stringify([...next]));
            return next;
        });
    }, []);

    const columnDefs: ColDef<Survivor>[] = columnConfigs.map(config => ({
        field: config.field,
        headerName: config.headerName,
        headerTooltip: config.headerTooltip,
        hide: !visibleColumns.has(config.field),
        lockVisible: config.locked,
        cellStyle: config.field === 'name' ? {textAlign: 'left'} : {textAlign: 'center'},
        headerClass: config.field === 'name' ? '' : 'ag-header-cell-center',
        flex: config.field === 'name' ? 2 : 1,
    }));

    return (
        <div className={styles.container} onContextMenu={(e) => e.preventDefault()}>
            <div className={styles.header}>
                <button
                    className={styles.btnGhost}
                    aria-label="Create Survivor"
                    title="Create Survivor"
                    onClick={onCreateSurvivor}
                >
                    <PlusIcon size={18} weight="bold" />
                </button>
                <div className={styles.columnMenuWrapper}>
                    <button
                        className="btn btn-ghost"
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
                                {columnConfigs.map(config => (
                                    <label
                                        key={config.field}
                                        className={`${styles.columnOption} ${config.locked ? styles.columnOptionDisabled : ''}`}
                                    >
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm"
                                            checked={visibleColumns.has(config.field)}
                                            onChange={() => toggleColumn(config.field)}
                                            disabled={config.locked}
                                        />
                                        <span className={styles.columnOptionLabel}>{config.headerName}</span>
                                    </label>
                                ))}
                            </div>
                        </>
                    )}
                </div>
            </div>
            <DataGrid
                rowData={data}
                defaultColDef={{cellStyle: {textAlign: 'center'}, headerClass: 'ag-header-cell-center'}}
                suppressCellFocus={true}
                columnDefs={columnDefs}
                autoSizeStrategy={{type: 'fitGridWidth'}}
                onCellContextMenu={handleCellContextMenu}
                onGridReady={(params) => {
                    gridApiRef.current = params.api;
                }}
                rowSelection={{mode: 'singleRow', hideDisabledCheckboxes: true, checkboxes: false}}
            />
            {contextMenu.visible && (
                <>
                    <div className={styles.overlay} onClick={closeMenu} onContextMenu={(e) => { e.preventDefault(); closeMenu(); }} />
                    <div
                        className={styles.contextMenu}
                        style={{left: contextMenu.x, top: contextMenu.y}}
                    >
                        <button
                            className={styles.contextMenuItem}
                            onClick={handleEdit}
                        >
                            Edit
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}