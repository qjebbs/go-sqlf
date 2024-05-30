package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

func (b *QueryBuilder) calcDependency() (map[TableAliased]bool, error) {
	tables := extractTables(
		b.selects,
		b.touches,
		b.conditions,
		b.orders,
		b.groupbys,
	)
	m := make(map[TableAliased]bool)
	// first table is the main table and always included
	m[b.tables[0]] = true
	for _, table := range tables {
		err := b.markDependencies(m, table)
		if err != nil {
			return nil, err
		}
	}
	// mark for CTEs
	for _, t := range b.tables {
		if b.distinct && b.froms[t].Optional && !m[t] {
			continue
		}
		// this could probably mark a CTE table that does not exists, but do no harm.
		m[NewTableAliased(t.Name, "")] = true
	}
	return m, nil
}

func (b *QueryBuilder) markDependencies(dep map[TableAliased]bool, t Table) error {
	ta, ok := b.appliedNames[t]
	if !ok {
		return fmt.Errorf("table not found: '%s'", t)
	}
	from, ok := b.froms[ta]
	if !ok {
		return fmt.Errorf("from undefined: '%s'", t)
	}
	if dep[ta] {
		return nil
	}
	dep[ta] = true
	for _, ft := range extractTables(from.Fragment) {
		if ft == t {
			continue
		}
		err := b.markDependencies(dep, ft)
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
