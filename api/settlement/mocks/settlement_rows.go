package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockRows struct {
	Rows    []SettlementRow
	cusrsor int
}

func (sr MockRows) Close() {

}
func (sr *MockRows) Err() error {
	panic("implement me")
}
func (sr *MockRows) CommandTag() pgconn.CommandTag {
	panic("implement me")
}
func (sr *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	panic("implement me")
}
func (sr *MockRows) Next() bool {
	return sr.cusrsor < len(sr.Rows)
}
func (sr *MockRows) Scan(dest ...any) error {
	sr.Rows[sr.cusrsor].Scan(dest...)
	sr.cusrsor++
	return nil
}
func (sr *MockRows) Values() ([]any, error) {
	panic("implement me")
}
func (sr *MockRows) RawValues() [][]byte {
	panic("implement me")
}
func (sr *MockRows) Conn() *pgx.Conn {
	panic("implement me")
}

type SettlementRow struct {
	Id                  int
	Owner               int
	Name                string
	SurvivalLimit       int
	DepartingSurvival   int
	CollectiveCognition int
	CurrentYear         int
}

func (s *SettlementRow) Scan(dest ...any) error {
	id := dest[0].(*int)
	owner := dest[1].(*int)
	name := dest[2].(*string)
	survivalLimit := dest[3].(*int)
	departingSurvival := dest[4].(*int)
	collectiveCognition := dest[5].(*int)
	currentYear := dest[6].(*int)

	*id = s.Id
	*owner = s.Owner
	*name = s.Name
	*survivalLimit = s.SurvivalLimit
	*departingSurvival = s.DepartingSurvival
	*collectiveCognition = s.CollectiveCognition
	*currentYear = s.CurrentYear

	return nil
}
