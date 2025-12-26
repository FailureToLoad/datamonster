import {useState} from 'react';
import {DataGrid} from '~/components/DataGrid';
import {type Survivor} from '~/types/survivor';
import type {CellContextMenuEvent} from 'ag-grid-community';

type CreateSurvivorProps = {
    data: Survivor[];
    onEditSurvivor?: (survivor: Survivor) => void;
};

type ContextMenuState = {
    visible: boolean;
    x: number;
    y: number;
    survivor: Survivor | null;
};

export function SurvivorTable({data, onEditSurvivor}: CreateSurvivorProps) {
    const [contextMenu, setContextMenu] = useState<ContextMenuState>({
        visible: false,
        x: 0,
        y: 0,
        survivor: null,
    });

    const closeMenu = () => setContextMenu(prev => ({...prev, visible: false}));

    const handleCellContextMenu = (event: CellContextMenuEvent<Survivor>) => {
        event.event?.preventDefault();
        const mouseEvent = event.event as MouseEvent;
        if (event.data) {
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
        setContextMenu(prev => ({...prev, visible: false}));
    };

    return (
        <div className="relative" onContextMenu={(e) => e.preventDefault()}>
            <DataGrid
                rowData={data}
                defaultColDef={{cellStyle: {textAlign: 'center'}, headerClass: 'ag-header-cell-center'}}
                suppressCellFocus={true}
                columnDefs={[
                    {field: 'name', cellStyle: {textAlign: 'left'}, headerClass: ''},
                    {field: 'gender'},
                    {field: 'birth'},
                    {field: 'huntxp'},
                    {field: 'survival'},
                    {field: 'movement'},
                    {field: 'speed'},
                    {field: 'strength'},
                    {field: 'accuracy'},
                    {field: 'evasion'},
                    {field: 'luck'},
                    {field: 'systemicpressure', headerName: 'S.P.', headerTooltip: 'System Pressure'},
                    {field: 'torment'},
                    {field: 'courage'},
                    {field: 'understanding'},
                    {field: 'insanity'},
                    {field: 'lumi'},
                ]}
                autoSizeStrategy={{type: 'fitGridWidth'}}
                onCellContextMenu={handleCellContextMenu}
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