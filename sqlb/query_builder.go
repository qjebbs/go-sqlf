package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

// QueryBuilder is the SQL query builder.
// It's recommended to wrap it with your struct to provide a
// more friendly API and improve fragment reusability.
type QueryBuilder struct {
	ctes     []*cte         // common table expressions in order
	ctesDict map[Table]*cte // the ctes by name, not alias

	tables     []*fromTable         // the tables in order
	tablesDict map[Table]*fromTable // the from tables by alias

	selects    *sqlf.Fragment         // select columns and keep values in scanning.
	touches    *sqlf.Fragment         // select columns but drop values in scanning.
	conditions *sqlf.Fragment         // where conditions, joined with AND.
	orders     *sqlf.Fragment         // order by columns, joined with comma.
	groupbys   *sqlf.Fragment         // group by columns, joined with comma.
	distinct   bool                   // select distinct
	limit      int64                  // limit count
	offset     int64                  // offset count
	unions     []sqlf.FragmentBuilder // union queries

	errors []error // errors during building

	debug bool // debug mode
}

type fromTable struct {
	Names    TableAliased
	Fragment *sqlf.Fragment
	Optional bool
}

// NewQueryBuilder returns a new QueryBuilder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		ctesDict:   make(map[Table]*cte),
		tablesDict: make(map[Table]*fromTable),
		selects:    sqlf.F("#join('#fragment', ', ')").WithPrefix("SELECT"),
		touches:    sqlf.F("#join('#fragment', ', ')"),
		conditions: sqlf.F("#join('#fragment', ' AND ')").WithPrefix("WHERE"),
		orders:     sqlf.F("#join('#fragment', ', ')").WithPrefix("ORDER BY"),
		groupbys:   sqlf.F("#join('#fragment', ', ')").WithPrefix("GROUP BY"),
	}
}

// Distinct set the flag for SELECT DISTINCT.
func (b *QueryBuilder) Distinct() *QueryBuilder {
	b.distinct = true
	return b
}

// Select replace the SELECT clause with the columns.
func (b *QueryBuilder) Select(columns ...*Column) *QueryBuilder {
	if len(columns) == 0 {
		return b
	}
	b.selects.WithFragments(convertFragmentBuilders(columns)...)
	return b
}

// Order is the sorting order.
type Order uint

// orders
const (
	Asc Order = iota
	AscNullsFirst
	AscNullsLast
	Desc
	DescNullsFirst
	DescNullsLast
)

var orders = []string{
	"ASC",
	"ASC NULLS FIRST",
	"ASC NULLS LAST",
	"DESC",
	"DESC NULLS FIRST",
	"DESC NULLS LAST",
}

// OrderBy set the sorting order. the order can be "ASC", "DESC", "ASC NULLS FIRST" or "DESC NULLS LAST"
func (b *QueryBuilder) OrderBy(column *Column, order Order) *QueryBuilder {
	idx := len(b.orders.Fragments) + 1
	alias := fmt.Sprintf("_order_%d", idx)

	if order > DescNullsLast {
		b.pushError(fmt.Errorf("invalid order: %d", order))
	}
	orderStr := orders[order]
	// pq: for SELECT DISTINCT, ORDER BY expressions must appear in select list
	b.touches.AppendFragments(sqlf.Ff("#f1 AS "+alias, column))
	b.orders.AppendFragments(sqlf.F(fmt.Sprintf("%s %s", alias, orderStr)))
	return b
}

// Limit set the limit.
func (b *QueryBuilder) Limit(limit int64) *QueryBuilder {
	if limit > 0 {
		b.limit = limit
	}
	return b
}

// Offset set the offset.
func (b *QueryBuilder) Offset(offset int64) *QueryBuilder {
	if offset > 0 {
		b.offset = offset
	}
	return b
}

// GroupBy set the sorting order.
func (b *QueryBuilder) GroupBy(columns ...*Column) *QueryBuilder {
	for _, c := range columns {
		b.groupbys.AppendFragments(c)
	}
	return b
}

// Union unions other query builders, the type of query builders can be
// *QueryBuilder or any other extended *QueryBuilder types (structs with
// *QueryBuilder embedded.)
func (b *QueryBuilder) Union(builders ...sqlf.FragmentBuilder) *QueryBuilder {
	b.unions = append(b.unions, builders...)
	return b
}
