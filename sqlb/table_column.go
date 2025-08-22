package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3"
)

var _ sqlf.Builder = (*Column)(nil)

// Column is a Column of a table.
// There are two ways to make columns:
//
//	foo := sqlb.NewTableAliased("foo", "f")
//	foo.Column("id") // "f.id"
//	sqlb.ExprColumn(sqlf.Ff("COALESCE(?.id,0)", foo)) // "COALESCE(f.id,0)"
type Column struct {
	fragment *sqlf.Fragment

	// the table of the anonymous column, set only for anonymous columns
	anonymousTable Table
}

// Build implements sqlf.Builder
func (c *Column) Build(ctx *sqlf.Context) (query string, err error) {
	// force build the table of the anonymous column,
	// so that its table dependency is tracked
	if c.anonymousTable != "" {
		_, err = c.anonymousTable.Build(ctx)
		if err != nil {
			return
		}
	}
	return c.fragment.Build(ctx)
}

// ExprColumn wraps a *Fragment of column expression to a *Column.
//
// A complex expression column is rather a fragment than a regular column,
// but the *Fragment cannot pass to methods requires the *Column type (like
// *QueryBuilder.Select()), this is where this function comes in.
//
//	t := sqlb.NewTableAliased("foo", "f")
//	sqlb.NewQueryBuilder().Select(
//		sqlb.ExprColumn(sqlf.Ff(
//			"COALESCE(?.id,0)",
//			t,
//		)),
//		sqlb.ExprColumn(sqlf.Ff(
//			"? > ? AS larger",
//			t.Column("bar"),
//			t.Column("baz"),
//		)),
//	)
//	// SELECT COALESCE(f.id,0), f.bar > f.baz AS larger ...
//
// Make sure all Tables are explicitly specified for the 'fragment' argument
// of sqlb.ExprColumn(), because they're required by *sqlb.QueryBuilder to
// calculate the table dependencies.
func ExprColumn(fragment *sqlf.Fragment) *Column {
	return &Column{
		fragment: fragment,
	}
}

// ExprColumns wraps a slice of *Fragment to a slice of *Column.
func ExprColumns(fragment ...*sqlf.Fragment) []*Column {
	r := make([]*Column, 0, len(fragment))
	for _, f := range fragment {
		r = append(r, ExprColumn(f))
	}
	return r
}
