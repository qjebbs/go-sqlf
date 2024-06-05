package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

func (b *QueryBuilder) collectDependencies() (map[TableAliased]bool, error) {
	tables := extractTables(
		b.selects,
		b.touches,
		b.conditions,
		b.orders,
		b.groupbys,
	)
	deps := make(map[TableAliased]bool)
	// first table is the main table and always included
	deps[b.tables[0].Names] = true
	for _, table := range tables {
		err := b.collectDepsFromTable(deps, table)
		if err != nil {
			return nil, err
		}
	}
	// mark for CTEs
	for _, t := range b.tables {
		if b.distinct && t.Optional && !deps[t.Names] {
			continue
		}
		if _, ok := b.ctesDict[t.Names.Name]; ok {
			deps[NewTableAliased(t.Names.Name, "")] = true
		}
	}
	return deps, nil
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
	for _, ft := range extractTables(from.Fragment) {
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

func extractTables(fragments ...sqlf.FragmentBuilder) []Table {
	tables := []Table{}
	dict := map[Table]bool{}
	extractTables2(fragments, &tables, dict)
	return tables
}

func extractTables2(fragments []sqlf.FragmentBuilder, tables *[]Table, dict map[Table]bool) {
	for _, f := range fragments {
		if f == nil {
			continue
		}
		if fragment, ok := f.(*sqlf.Fragment); ok {
			extractTables2(fragment.Fragments, tables, dict)
			continue
		}
		if column, ok := f.(*Column); ok && column != nil {
			if column.table != "" {
				if !dict[column.table] {
					collectTable(column.table, tables, dict)
				}
			} else {
				extractTables2(column.fragment.Fragments, tables, dict)
			}
			continue
		}

		if table, ok := f.(Table); ok {
			collectTable(table, tables, dict)
			continue
		}

		if table, ok := f.(TableAliased); ok {
			collectTable(table.AppliedName(), tables, dict)
		}
	}
}

func collectTable(t Table, tables *[]Table, dict map[Table]bool) {
	if dict[t] {
		return
	}
	*tables = append(*tables, t)
	dict[t] = true
}
