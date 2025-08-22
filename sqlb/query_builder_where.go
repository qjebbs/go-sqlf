package sqlb

import (
	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/util"
)

// Where add a condition.  e.g.:
//
//	b.Where(
//		sqlf.F("? = ?", a.Column("id"), 1),
//	)
func (b *QueryBuilder) Where(s *sqlf.Fragment) *QueryBuilder {
	if s == nil {
		return b
	}
	b.conditions = append(b.conditions, s)
	return b
}

// Where2 is a helper func similar to Where(), which adds a simple where condition. e.g.:
//
//	b.Where2(column, "=", 1)
//
// it's equivalent to:
//
//	b.Where(
//		sqlf.F("? = ?", column, 1),
//	)
func (b *QueryBuilder) Where2(column *Column, op string, arg any) *QueryBuilder {
	b.conditions = append(
		b.conditions,
		sqlf.F("?"+op+"?", column, arg),
	)
	return b
}

// WhereIn adds a where IN condition like `t.id IN (1,2,3)`
func (b *QueryBuilder) WhereIn(column *Column, list any) *QueryBuilder {
	return b.Where(
		sqlf.F(
			"? IN (?)",
			column,
			sqlf.Join(", ", util.ArgsFlatted(list)...),
		),
	)
}

// WhereNotIn adds a where NOT IN condition like `t.id NOT IN (1,2,3)`
func (b *QueryBuilder) WhereNotIn(column *Column, list any) *QueryBuilder {
	return b.Where(
		sqlf.F(
			"? NOT IN (?)",
			column,
			sqlf.Join(", ", util.ArgsFlatted(list)...),
		),
	)
}
