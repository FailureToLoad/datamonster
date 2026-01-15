import { describe, it, expect } from 'vitest';
import { render, screen } from '@solidjs/testing-library';
import userEvent from '@testing-library/user-event';
import { DataTable, type ColumnConfig } from './DataTable';

type TestRow = {
  id: string;
  name: string;
  age: number;
  email: string;
  status: string;
};

const testColumns: ColumnConfig<TestRow>[] = [
  { field: 'name', headerName: 'Name', locked: true, default: true, sortable: true },
  { field: 'age', headerName: 'Age', locked: false, default: true, sortable: true },
  { field: 'email', headerName: 'Email', locked: false, default: true, sortable: true },
  { field: 'status', headerName: 'Status', locked: false, default: false, sortable: true },
];

const testRows: TestRow[] = [
  { id: '1', name: 'Alice', age: 30, email: 'alice@example.com', status: 'active' },
  { id: '2', name: 'Bob', age: 25, email: 'bob@example.com', status: 'inactive' },
  { id: '3', name: 'Charlie', age: 35, email: 'charlie@example.com', status: 'active' },
  { id: '4', name: 'Diana', age: 28, email: 'diana@example.com', status: 'pending' },
];

function getColumnHeaders(): string[] {
  const headers = screen.getAllByRole('columnheader');
  return headers.map((th) => th.textContent?.replace(/[▲▼]/g, '').trim() ?? '');
}

function getColumnCells(columnIndex: number): string[] {
  const rows = screen.getAllByRole('row').slice(1);
  return rows.map((row) => {
    const cells = row.querySelectorAll('td');
    return cells[columnIndex]?.textContent ?? '';
  });
}

describe('DataTable', () => {
  describe('Column Visibility', () => {
    it('shows default columns on initial render', () => {
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const headers = getColumnHeaders();
      expect(headers).toContain('Name');
      expect(headers).toContain('Age');
      expect(headers).toContain('Email');
      expect(headers).not.toContain('Status');
    });

    it('can add a hidden column via column menu', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      expect(getColumnHeaders()).not.toContain('Status');

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);

      const statusCheckbox = screen.getByRole('checkbox', { name: /status/i });
      await user.click(statusCheckbox);

      expect(getColumnHeaders()).toContain('Status');
    });

    it('can remove a non-locked column via column menu', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      expect(getColumnHeaders()).toContain('Age');

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);

      const ageCheckbox = screen.getByRole('checkbox', { name: /age/i });
      await user.click(ageCheckbox);

      expect(getColumnHeaders()).not.toContain('Age');
    });

    it('does not show locked columns in the column menu', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);

      const checkboxes = screen.getAllByRole('checkbox');
      const labels = checkboxes.map((cb) => cb.closest('label')?.textContent?.trim());

      expect(labels).not.toContain('Name');
      expect(labels).toContain('Age');
      expect(labels).toContain('Email');
      expect(labels).toContain('Status');
    });

    it('locked name column is always visible', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      expect(getColumnHeaders()).toContain('Name');

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);

      const nameCheckbox = screen.queryByRole('checkbox', { name: /^name$/i });
      expect(nameCheckbox).toBeNull();

      expect(getColumnHeaders()).toContain('Name');
    });

    it('can toggle columns on and off multiple times', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const settingsButton = screen.getByTitle('Configure columns');

      // Hide Email column
      await user.click(settingsButton);
      const emailCheckbox = screen.getByRole('checkbox', { name: /email/i });
      await user.click(emailCheckbox);
      expect(getColumnHeaders()).not.toContain('Email');

      // Show Email column again (menu is still open, checkbox is still there)
      const emailCheckboxAgain = screen.getByRole('checkbox', { name: /email/i });
      await user.click(emailCheckboxAgain);
      expect(getColumnHeaders()).toContain('Email');
    });
  });

  describe('Sorting', () => {
    it('sorts name column ascending on first click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      await user.click(nameHeader);

      const names = getColumnCells(0);
      expect(names).toEqual(['Alice', 'Bob', 'Charlie', 'Diana']);
    });

    it('sorts name column descending on second click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      await user.click(nameHeader);
      await user.click(nameHeader);

      const names = getColumnCells(0);
      expect(names).toEqual(['Diana', 'Charlie', 'Bob', 'Alice']);
    });

    it('clears sorting on third click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      await user.click(nameHeader);
      await user.click(nameHeader);
      await user.click(nameHeader);

      const names = getColumnCells(0);
      expect(names).toEqual(['Alice', 'Bob', 'Charlie', 'Diana']);
    });

    it('sorts age column descending on first click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const ageHeader = screen.getByRole('columnheader', { name: /age/i });
      await user.click(ageHeader);

      const names = getColumnCells(0);
      // TanStack Table sorts numeric columns descending first by default
      // Charlie(35), Alice(30), Diana(28), Bob(25)
      expect(names).toEqual(['Charlie', 'Alice', 'Diana', 'Bob']);
    });

    it('sorts age column ascending on second click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const ageHeader = screen.getByRole('columnheader', { name: /age/i });
      await user.click(ageHeader);
      await user.click(ageHeader);

      const names = getColumnCells(0);
      // Bob(25), Diana(28), Alice(30), Charlie(35)
      expect(names).toEqual(['Bob', 'Diana', 'Alice', 'Charlie']);
    });

    it('sorts email column ascending on first click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const emailHeader = screen.getByRole('columnheader', { name: /email/i });
      await user.click(emailHeader);

      const emails = getColumnCells(2);
      expect(emails).toEqual([
        'alice@example.com',
        'bob@example.com',
        'charlie@example.com',
        'diana@example.com',
      ]);
    });

    it('sorts email column descending on second click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const emailHeader = screen.getByRole('columnheader', { name: /email/i });
      await user.click(emailHeader);
      await user.click(emailHeader);

      const emails = getColumnCells(2);
      expect(emails).toEqual([
        'diana@example.com',
        'charlie@example.com',
        'bob@example.com',
        'alice@example.com',
      ]);
    });

    it('sorts status column ascending on first click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);
      const statusCheckbox = screen.getByRole('checkbox', { name: /status/i });
      await user.click(statusCheckbox);

      const statusHeader = screen.getByRole('columnheader', { name: /status/i });
      await user.click(statusHeader);

      const statuses = getColumnCells(3);
      expect(statuses).toEqual(['active', 'active', 'inactive', 'pending']);
    });

    it('sorts status column descending on second click', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const settingsButton = screen.getByTitle('Configure columns');
      await user.click(settingsButton);
      const statusCheckbox = screen.getByRole('checkbox', { name: /status/i });
      await user.click(statusCheckbox);

      const statusHeader = screen.getByRole('columnheader', { name: /status/i });
      await user.click(statusHeader);
      await user.click(statusHeader);

      const statuses = getColumnCells(3);
      expect(statuses).toEqual(['pending', 'inactive', 'active', 'active']);
    });

    it('shows ascending indicator when sorted ascending', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      await user.click(nameHeader);

      expect(nameHeader.textContent).toContain('▲');
    });

    it('shows descending indicator when sorted descending', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      await user.click(nameHeader);
      await user.click(nameHeader);

      expect(nameHeader.textContent).toContain('▼');
    });

    it('switching sort column clears previous sort', async () => {
      const user = userEvent.setup();
      render(() => <DataTable columns={testColumns} rows={testRows} />);

      const nameHeader = screen.getByRole('columnheader', { name: /name/i });
      const ageHeader = screen.getByRole('columnheader', { name: /age/i });

      // Sort name descending
      await user.click(nameHeader);
      await user.click(nameHeader);
      expect(nameHeader.textContent).toContain('▼');

      // Switch to age column - numeric columns sort descending first
      await user.click(ageHeader);
      expect(ageHeader.textContent).toContain('▼');

      // Verify rows are sorted by age descending: Charlie(35), Alice(30), Diana(28), Bob(25)
      const names = getColumnCells(0);
      expect(names).toEqual(['Charlie', 'Alice', 'Diana', 'Bob']);
    });
  });
});
