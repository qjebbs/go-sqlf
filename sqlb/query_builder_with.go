package sqlb

import "github.com/qjebbs/go-sqlf"

// With adds a fragment as common table expression, the built query of s should be a subquery.
func (b *QueryBuilder) With(name sqlf.Table, builder sqlf.Builder) *QueryBuilder {
	b.ctes = append(b.ctes, &cte{
		table:   NewTable(name, ""),
		Builder: builder,
	})
	return b
}

type cte struct {
	table Table
	sqlf.Builder
}
