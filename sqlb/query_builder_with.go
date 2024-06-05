package sqlb

import "github.com/qjebbs/go-sqlf/v2"

// Why it's impossible to colloect dependencies between CTEs?
//
// Consider the following query:
//
//	sqlb.NewQueryBuilder().
//		With(cteA, ...).
//		With(cteB, ...).
//		Select(cteB.Column("*")).
//		From(cteB)
//
// We cannot determine without semantic analysis,
// when cteB references with a same identifier as cteA, we cannot tell
// if cteA is self-contained by cteB, or cteB requires cteA, without
// semantic analysis. Not to mention that cteB could be any type
// implementing FragmentBuilder, and it's impossible to determine.

// With adds a fragment as common table expression,
// the built query of s should be a subquery,
// deps are the other CTEs that the CTE depends on.
//
// CTE dependencies are not automatically calculated, since it's
// not possible to do so without semantic analysis.
func (b *QueryBuilder) With(name Table, builder sqlf.FragmentBuilder, deps ...Table) *QueryBuilder {
	cte := &cte{
		name:            name,
		deps:            deps,
		FragmentBuilder: builder,
	}
	b.ctes = append(b.ctes, cte)
	b.ctesDict[name] = cte
	return b
}

type cte struct {
	name Table
	deps []Table
	sqlf.FragmentBuilder
}
