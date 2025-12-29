import {useState, useRef, useCallback} from 'react';
import {DataGrid} from '~/components/DataGrid';
import {type Survivor} from '~/types/survivor';
import type {CellContextMenuEvent, GridApi, ColDef} from 'ag-grid-community';
import {GearIcon} from '@phosphor-icons/react';
import NewSurvivorDialog from './survivorDialog';

const COLUMN_CONFIG_KEY = 'survivor-table-columns';

type SurvivorTableProps = {
    data: Survivor[];
    settlementId: string;
    onEditSurvivor?: (survivor: Survivor) => void;
    onSurvivorCreated?: () => void;
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

export function SurvivorTable({data, settlementId, onEditSurvivor, onSurvivorCreated}: SurvivorTableProps) {
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
        <div className="relative" onContextMenu={(e) => e.preventDefault()}>
            <div className="flex justify-between items-center mb-2">
                <NewSurvivorDialog
                    settlementId={settlementId}
                    onSuccess={onSurvivorCreated}
                />
                <div className="relative">
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
                                className="fixed inset-0 z-40"
                                onClick={() => setColumnMenuOpen(false)}
                            />
                            <div className="absolute right-0 top-full z-50 bg-base-100 border border-base-300 rounded-lg shadow-lg py-2 min-w-48">
                                <p className="px-3 py-1 text-xs font-semibold text-base-content/60 uppercase">Columns</p>
                                {columnConfigs.map(config => (
                                    <label
                                        key={config.field}
                                        className={`flex items-center gap-2 px-3 py-1.5 hover:bg-base-200 cursor-pointer ${config.locked ? 'opacity-50 cursor-not-allowed' : ''}`}
                                    >
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm"
                                            checked={visibleColumns.has(config.field)}
                                            onChange={() => toggleColumn(config.field)}
                                            disabled={config.locked}
                                        />
                                        <span className="text-sm">{config.headerName}</span>
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
                    <div className="fixed inset-0 z-40" onClick={closeMenu} onContextMenu={(e) => { e.preventDefault(); closeMenu(); }} />
                    <div
                        className="fixed z-50 bg-base-100 border border-base-300 rounded-lg shadow-lg py-1 min-w-32"
                        style={{left: contextMenu.x, top: contextMenu.y}}
                    >
                        <button
                            className="w-full px-4 py-2 text-left hover:bg-base-200 transition-colors"
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