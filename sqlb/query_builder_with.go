package sqlb

import "github.com/qjebbs/go-sqlf/v2"

// With adds a fragment as common table expression, the built query of s should be a subquery.
func (b *QueryBuilder) With(name Table, builder sqlf.FragmentBuilder) *QueryBuilder {
	b.ctes = append(b.ctes, &cte{
		table:           NewTableAliased(name, ""),
		FragmentBuilder: builder,
	})
	return b
}

type cte struct {
	table TableAliased
	sqlf.FragmentBuilder
}
