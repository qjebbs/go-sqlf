package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/syntax"
)

var _ sqlf.Builder = (*QueryBuilder)(nil)

// QueryBuilder is the SQL query builder.
// It's recommended to wrap it with your struct to provide a
// more friendly API and improve fragment reusability.
type QueryBuilder struct {
	bindVarStyle syntax.BindVarStyle // the bindvar style

	ctes         []*cte               // common table expressions
	froms        map[Table]*fromTable // the from tables by alias
	tables       []Table              // the tables in order
	appliedNames map[sqlf.Table]Table // applied table name mapping, the name is alias, or name if alias is empty

	selects    *sqlf.Fragment // select columns and keep values in scanning.
	touches    *sqlf.Fragment // select columns but drop values in scanning.
	conditions *sqlf.Fragment // where conditions, joined with AND.
	orders     *sqlf.Fragment // order by columns, joined with comma.
	groupbys   *sqlf.Fragment // group by columns, joined with comma.
	distinct   bool           // select distinct
	limit      int64          // limit count
	offset     int64          // offset count
	unions     []sqlf.Builder // union queries

	errors []error // errors during building

	debug bool // debug mode
}

type fromTable struct {
	Fragment *sqlf.Fragment
	Optional bool
}

// NewQueryBuilder returns a new QueryBuilder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		froms:        map[Table]*fromTable{},
		appliedNames: make(map[sqlf.Table]Table),
		selects: &sqlf.Fragment{
			Prefix: "SELECT",
			Raw:    "#join('#column', ', ')",
		},
		touches: &sqlf.Fragment{
			Prefix: "",
			Raw:    "#join('#fragment', ', ')",
		},
		conditions: &sqlf.Fragment{
			Prefix: "WHERE",
			Raw:    "#join('#fragment', ' AND ')",
		},
		orders: &sqlf.Fragment{
			Prefix: "ORDER BY",
			Raw:    "#join('#fragment', ', ')",
		},
		groupbys: &sqlf.Fragment{
			Prefix: "GROUP BY",
			Raw:    "#join('#fragment', ', ')",
		},
	}
}

// Distinct set the flag for SELECT DISTINCT.
func (b *QueryBuilder) Distinct() *QueryBuilder {
	b.distinct = true
	return b
}

// Select replace the SELECT clause with the columns.
func (b *QueryBuilder) Select(columns ...*sqlf.TableColumn) *QueryBuilder {
	if len(columns) == 0 {
		return b
	}
	b.selects.WithColumns(columns...)
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
func (b *QueryBuilder) OrderBy(column *sqlf.TableColumn, order Order) *QueryBuilder {
	idx := len(b.orders.Fragments) + 1
	alias := fmt.Sprintf("_order_%d", idx)

	if order > DescNullsLast {
		b.pushError(fmt.Errorf("invalid order: %d", order))
	}
	orderStr := orders[order]
	// pq: for SELECT DISTINCT, ORDER BY expressions must appear in select list
	b.touches.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1 AS " + alias,
		Columns: []*sqlf.TableColumn{column},
	})
	b.orders.AppendFragments(&sqlf.Fragment{
		Raw:     fmt.Sprintf("%s %s", alias, orderStr),
		Columns: nil,
		Args:    nil,
	})
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
func (b *QueryBuilder) GroupBy(column *sqlf.TableColumn, args ...any) *QueryBuilder {
	b.groupbys.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1",
		Columns: []*sqlf.TableColumn{column},
		Args:    args,
	})
	return b
}

// Union unions other query builders, the type of query builders can be
// *QueryBuilder or any other extended *QueryBuilder types (structs with
// *QueryBuilder embedded.)
func (b *QueryBuilder) Union(builders ...sqlf.Builder) *QueryBuilder {
	b.unions = append(b.unions, builders...)
	return b
}

// BindVar set the bindvar style.
func (b *QueryBuilder) BindVar(style syntax.BindVarStyle) *QueryBuilder {
	b.bindVarStyle = style
	return b
}
