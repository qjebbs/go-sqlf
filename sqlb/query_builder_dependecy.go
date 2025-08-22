package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/syntax"
	"github.com/qjebbs/go-sqlf/v3/util"
)

type depTablesKey struct{}

func contextWithDeps(ctx *sqlf.Context, deps map[Table]bool) *sqlf.Context {
	ctx.WithValue(depTablesKey{}, deps)
	return ctx
}

func depsFromContext(ctx *sqlf.Context) map[Table]bool {
	dep := ctx.Value(depTablesKey{})
	if dep == nil {
		return nil
	}
	return dep.(map[Table]bool)
}

// collectDependencies collects the dependencies of the tables.
func (b *QueryBuilder) collectDependencies() (map[TableAliased]bool, error) {
	builders := util.ArgsFlatted(
		b.selects,
		b.touches,
		b.conditions,
		b.groupbys,
		b.havings,
	)
	for _, order := range b.orders {
		builders = append(builders, order.column)
	}

	tables, err := extractTables(builders)
	if err != nil {
		return nil, fmt.Errorf("collect dependencies: %w", err)
	}
	deps := make(map[TableAliased]bool)
	// first table is the main table and always included
	deps[b.tables[0].Names] = true
	for table := range tables {
		err := b.collectDepsFromTable(deps, table)
		if err != nil {
			return nil, err
		}
	}
	// mark for CTEs
	depsCTE := make(map[TableAliased]bool)
	for _, t := range b.tables {
		if (b.distinct || len(b.groupbys) > 0) && t.Optional && !deps[t.Names] {
			continue
		}
		if cte, ok := b.ctesDict[t.Names.AppliedName()]; ok {
			b.collectDepsFromCTE(depsCTE, cte)
		}
	}
	for cte := range depsCTE {
		deps[cte] = true
	}
	return deps, nil
}

func (b *QueryBuilder) collectDepsFromCTE(deps map[TableAliased]bool, cte *cte) error {
	key := cte.name
	if deps[key] {
		return nil
	}
	deps[key] = true
	tables, err := extractTables([]any{cte.Builder})
	if err != nil {
		return fmt.Errorf("collect dependencies of CTE %q: %w", cte.name, err)
	}
	for dep := range tables {
		if cte, ok := b.ctesDict[dep]; ok {
			err := b.collectDepsFromCTE(deps, cte)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *QueryBuilder) collectDepsFromTable(dep map[TableAliased]bool, t Table) error {
	from, ok := b.tablesDict[t]
	if !ok {
		return fmt.Errorf("from undefined: '%s'", t)
	}
	if dep[from.Names] {
		return nil
	}
	dep[from.Names] = true
	tables, err := extractTables(from.Fragment.Args)
	if err != nil {
		return fmt.Errorf("collect dependencies of table %q: %w", from.Names.Name, err)
	}
	for ft := range tables {
		if ft == t {
			continue
		}
		err := b.collectDepsFromTable(dep, ft)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractTables(args []any) (map[Table]bool, error) {
	tables := make(map[Table]bool)
	ctx := contextWithDeps(sqlf.NewContext(syntax.Dollar), tables)
	_, err := sqlf.Join(";", args...).Build(ctx)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
