package sqlb

import (
	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/util"
)

// Where add a condition.  e.g.:
//
//	b.Where(
//		sqlf.F("#f1 = $1").
//			WithFragments(a.Column("id")).
//			WithArgs(1),
//	)
func (b *QueryBuilder) Where(s *sqlf.Fragment) *QueryBuilder {
	if s == nil {
		return b
	}
	b.conditions.AppendArgs(s)
	return b
}

// Where2 is a helper func similar to Where(), which adds a simple where condition. e.g.:
//
//	b.Where2(column, "=", 1)
//
// it's equivalent to:
//
//	b.Where(
//		sqlf.F("#f1 = $1").
//			WithFragments(column).
//			WithArgs(1),
//	)
func (b *QueryBuilder) Where2(column *Column, op string, arg any) *QueryBuilder {
	b.conditions.AppendArgs(
		sqlf.F("#f1" + op + "$1").
			AppendArgs(column).
			WithArgs(arg),
	)
	return b
}

// WhereIn adds a where IN condition like `t.id IN (1,2,3)`
func (b *QueryBuilder) WhereIn(column *Column, list any) *QueryBuilder {
	return b.Where(
		sqlf.F("#f1 IN (#join('#arg', ', '))").
			WithArgs(column).
			AppendArgs(util.ArgsFlatted(list)...),
	)
}

// WhereNotIn adds a where NOT IN condition like `t.id NOT IN (1,2,3)`
func (b *QueryBuilder) WhereNotIn(column *Column, list any) *QueryBuilder {
	return b.Where(
		sqlf.F("#f1 NOT IN (#join('#arg', ', '))").
			WithArgs(column).
			AppendArgs(util.ArgsFlatted(list)...),
	)
}
