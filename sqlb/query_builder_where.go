package sqlb

import (
	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/util"
)

// Where add a condition.  e.g.:
//
//	b.Where(&sqlf.Fragment{
//		Raw: "#c1 = $1",
//		Columns: t.Columns("id"),
//		Args: []any{1},
//	})
func (b *QueryBuilder) Where(s *sqlf.Fragment) *QueryBuilder {
	if s == nil {
		return b
	}
	b.conditions.AppendFragments(s)
	return b
}

// Where2 is a helper func similar to Where(), which adds a simple where condition. e.g.:
//
//	b.Where2(column, "=", 1)
//
// it's  equivalent to:
//
//	b.Where(&sqlf.Fragment{
//		Raw: "#c1=$1",
//		Columns: []Column{column},
//		Args: []any{1},
//	})
func (b *QueryBuilder) Where2(column *sqlf.TableColumn, op string, arg any) *QueryBuilder {
	b.conditions.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1" + op + "$1",
		Columns: []*sqlf.TableColumn{column},
		Args:    []any{arg},
	})
	return b
}

// WhereIn adds a where IN condition like `t.id IN (1,2,3)`
func (b *QueryBuilder) WhereIn(column *sqlf.TableColumn, list any) *QueryBuilder {
	return b.Where(&sqlf.Fragment{
		Raw:     "#c1 IN (#join('#$', ', '))",
		Columns: []*sqlf.TableColumn{column},
		Args:    util.Args(list),
	})
}

// WhereNotIn adds a where NOT IN condition like `t.id NOT IN (1,2,3)`
func (b *QueryBuilder) WhereNotIn(column *sqlf.TableColumn, list any) *QueryBuilder {
	return b.Where(&sqlf.Fragment{
		Raw:     "#c1 NOT IN (#join('#$', ', '))",
		Columns: []*sqlf.TableColumn{column},
		Args:    util.Args(list),
	})
}
