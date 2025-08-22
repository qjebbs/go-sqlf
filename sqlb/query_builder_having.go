package sqlb

import (
	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/util"
)

// Having add a condition.  e.g.:
//
//	b.Having(
//		sqlf.F("#f1 = $1").
//			WithFragments(a.Column("id")).
//			WithArgs(1),
//	)
func (b *QueryBuilder) Having(s *sqlf.Fragment) *QueryBuilder {
	if s == nil {
		return b
	}
	b.havings.AppendArgs(s)
	return b
}

// Having2 is a helper func similar to Having(), which adds a simple where condition. e.g.:
//
//	b.Having2(column, "=", 1)
//
// it's equivalent to:
//
//	b.Having(
//		sqlf.F("#f1 = $1").
//			WithFragments(column).
//			WithArgs(1),
//	)
func (b *QueryBuilder) Having2(column *Column, op string, arg any) *QueryBuilder {
	b.havings.AppendArgs(
		sqlf.F("#f1" + op + "$1").
			WithArgs(column).
			AppendArgs(arg),
	)
	return b
}

// HavingIn adds a where IN condition like `t.id IN (1,2,3)`
func (b *QueryBuilder) HavingIn(column *Column, list any) *QueryBuilder {
	return b.Having(
		sqlf.F("#f1 IN (#join('#arg', ', '))").
			WithArgs(column).
			AppendArgs(util.ArgsFlatted(list)...),
	)
}

// HavingNotIn adds a where NOT IN condition like `t.id NOT IN (1,2,3)`
func (b *QueryBuilder) HavingNotIn(column *Column, list any) *QueryBuilder {
	return b.Having(
		sqlf.F("#f1 NOT IN (#join('#arg', ', '))").
			WithArgs(column).
			AppendArgs(util.ArgsFlatted(list)...),
	)
}
