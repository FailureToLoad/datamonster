package mocks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockDatabase struct {
	Rows *MockRows
	Row  *SettlementRow
}

func (sd MockDatabase) Close() {
	fmt.Println("Close called")
}
func (sd MockDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("not implemented")
}
func (sd MockDatabase) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	panic("not implemented")
}
func (sd MockDatabase) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error) {
	if sd.Rows == nil {
		panic("rows field not set")
	}
	return sd.Rows, nil
}
func (sd MockDatabase) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
	if sd.Row == nil {
		panic("row field not set")
	}

	return sd.Row
}

func (sd *MockDatabase) SetRows(rows *MockRows) {
	sd.Rows = rows
}

func (sd *MockDatabase) SetRow(row *SettlementRow) {
	sd.Row = row
}
