package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf"
)

func (b *QueryBuilder) calcDependency() (map[Table]bool, error) {
	tables := extractTables(
		b.selects,
		b.touches,
		b.conditions,
		b.orders,
		b.groupbys,
	)
	m := make(map[Table]bool)
	// first table is the main table and always included
	m[b.tables[0]] = true
	for _, t := range tables {
		err := b.markDependencies(m, t.Table)
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
		m[NewTable(t.Name, "")] = true
	}
	return m, nil
}

func (b *QueryBuilder) markDependencies(dep map[Table]bool, t sqlf.Table) error {
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
		if ft.Table == t {
			continue
		}
		err := b.markDependencies(dep, ft.Table)
		if err != nil {
			return fmt.Errorf("%s: %s", ft.Source, err)
		}
	}
	return nil
}

type tableWithSouce struct {
	Table  sqlf.Table
	Source string
}

func extractTables(fragments ...*sqlf.Fragment) []*tableWithSouce {
	tables := []*tableWithSouce{}
	dict := map[sqlf.Table]bool{}
	extractTables2(fragments, &tables, &dict)
	return tables
}

func extractTables2(fragments []*sqlf.Fragment, tables *[]*tableWithSouce, dict *map[sqlf.Table]bool) {
	for _, f := range fragments {
		if f == nil {
			continue
		}
		for i, t := range f.Tables {
			if (*dict)[t] {
				continue
			}
			*tables = append(*tables, &tableWithSouce{
				Table:  t,
				Source: fmt.Sprintf("#tables%d of '%s'", i+1, f.Raw),
			})
			(*dict)[t] = true
		}
		for i, c := range f.Columns {
			if c == nil || (*dict)[c.Table] {
				continue
			}
			*tables = append(*tables, &tableWithSouce{
				Table:  c.Table,
				Source: fmt.Sprintf("#column%d '%s' of '%s'", i+1, c.Raw, f.Raw),
			})
			(*dict)[c.Table] = true
		}
		extractTables2(f.Fragments, tables, dict)
	}
}
