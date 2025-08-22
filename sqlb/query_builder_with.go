package sqlb

import "github.com/qjebbs/go-sqlf/v3"

// With adds a fragment as common table expression,
// the built query of s should be a subquery,
// deps are the other CTEs that the CTE depends on.
//
// CTE dependencies are not automatically calculated, since it's
// not possible to do so without semantic analysis.
func (b *QueryBuilder) With(name TableAliased, builder sqlf.Builder) *QueryBuilder {
	cte := &cte{
		name:    name,
		Builder: builder,
	}
	b.ctes = append(b.ctes, cte)
	b.ctesDict[name.AppliedName()] = cte
	return b
}

type cte struct {
	name TableAliased
	sqlf.Builder
}
