package querybuilder

import (
	"fmt"
	"strings"

	"github.com/failuretoload/datamonster/pkg/querybuilder/columns"
)

type Table interface {
	TableName() string
	AllColumns() []string
}

type Query interface {
	Build() (string, []interface{})
}

type WhereClause struct {
	condition columns.Condition
}

func Where(condition columns.Condition) WhereClause {
	return WhereClause{condition: condition}
}

func (w WhereClause) Build() (string, []interface{}) {
	if w.condition == nil {
		return "", nil
	}

	whereClause, args, _ := w.condition.Build(1)
	if whereClause == "" {
		return "", nil
	}

	return "WHERE " + whereClause, args
}

type SelectQuery struct {
	table      Table
	selections []string
	where      *WhereClause
	orderBy    []string
	limit      *int
	offset     *int
}

func Select(table Table, cols ...string) *SelectQuery {
	if len(cols) == 0 {
		cols = table.AllColumns()
	}
	return &SelectQuery{
		table:      table,
		selections: cols,
	}
}

func (q *SelectQuery) Where(w WhereClause) *SelectQuery {
	q.where = &w
	return q
}

func (q *SelectQuery) OrderBy(cols ...string) *SelectQuery {
	q.orderBy = cols
	return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.limit = &limit
	return q
}

func (q *SelectQuery) Offset(offset int) *SelectQuery {
	q.offset = &offset
	return q
}

func (q *SelectQuery) Build() (string, []interface{}) {
	var params []interface{}
	query := &strings.Builder{}

	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.selections, ", "))
	query.WriteString(" FROM ")
	query.WriteString(q.table.TableName())

	if q.where != nil {
		if whereClause, whereArgs := q.where.Build(); whereClause != "" {
			query.WriteString(" ")
			query.WriteString(whereClause)
			params = append(params, whereArgs...)
		}
	}

	if len(q.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(q.orderBy, ", "))
	}

	if q.limit != nil {
		query.WriteString(" LIMIT ")
		query.WriteString(fmt.Sprint(*q.limit))
	}

	if q.offset != nil {
		query.WriteString(" OFFSET ")
		query.WriteString(fmt.Sprint(*q.offset))
	}

	return query.String(), params
}
