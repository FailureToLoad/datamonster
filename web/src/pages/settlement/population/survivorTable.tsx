import { DataTable, type ColumnConfig } from '~/components/DataTable';
import { type Survivor } from '~/lib/survivor';

type SurvivorTableProps = {
    data: Survivor[];
    onEditSurvivor?: (survivor: Survivor) => void;
    onCreateSurvivor?: () => void;
};

const columns: ColumnConfig<Survivor>[] = [
    { field: 'name', headerName: 'Name', locked: true, default: true, sortable: true },
    { field: 'gender', headerName: 'Gender', locked: false, default: true },
    { field: 'birth', headerName: 'Birth', locked: false },
    { field: 'huntxp', headerName: 'XP', headerTooltip: 'Hunt XP', locked: false, default: true, sortable: true },
    { field: 'survival', headerName: 'Survival', locked: false, default: true, sortable: true },
    { field: 'movement', headerName: 'Movement', locked: false, sortable: true },
    { field: 'speed', headerName: 'Speed', locked: false, sortable: true },
    { field: 'strength', headerName: 'Strength', locked: false, sortable: true },
    { field: 'accuracy', headerName: 'Accuracy', locked: false, sortable: true },
    { field: 'evasion', headerName: 'Evasion', locked: false, sortable: true },
    { field: 'luck', headerName: 'Luck', locked: false, sortable: true },
    { field: 'systemicPressure', headerName: 'S.P.', headerTooltip: 'Systemic Pressure', locked: false, sortable: true },
    { field: 'torment', headerName: 'Torment', locked: false, sortable: true },
    { field: 'courage', headerName: 'Courage', locked: false, default: true, sortable: true },
    { field: 'understanding', headerName: 'Understanding', locked: false, default: true, sortable: true },
    { field: 'insanity', headerName: 'Insanity', locked: false, default: true, sortable: true },
    { field: 'lumi', headerName: 'Lumi', locked: false, sortable: true },
];

export function SurvivorTable({ data, onEditSurvivor, onCreateSurvivor }: SurvivorTableProps) {
    return (
        <DataTable
            columns={columns}
            rows={data}
            contextActions={onEditSurvivor && [{ label: 'Edit', onClick: onEditSurvivor }]}
            addRow={onCreateSurvivor && { label: 'Create Survivor', onClick: onCreateSurvivor }}
            emptyMessage="No Survivors"
        />
    );
}
