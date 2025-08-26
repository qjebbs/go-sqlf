package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v3"
)

// From set the from table.
func (b *QueryBuilder) From(t TableAliased) *QueryBuilder {
	if t.Name == "" {
		b.pushError(fmt.Errorf("from table is empty"))
		return b
	}
	table := &fromTable{
		Names:    t,
		Builder:  NewTableAsBuilder(t),
		Optional: false,
	}
	if len(b.tables) == 0 {
		b.tables = append(b.tables, table)
	} else {
		b.tables[0] = table
	}
	b.tablesDict[t.AppliedName()] = table
	return b
}

// InnerJoin append a inner join table.
func (b *QueryBuilder) InnerJoin(t TableAliased, on *sqlf.Fragment) *QueryBuilder {
	return b.join("INNER JOIN", t, on, false)
}

// LeftJoin append / replace a left join table.
func (b *QueryBuilder) LeftJoin(t TableAliased, on *sqlf.Fragment) *QueryBuilder {
	return b.join("LEFT JOIN", t, on, false)
}

// LeftJoinOptional appends or replaces a LEFT JOIN table and marks it as optional.
// The optional join will be removed if no columns from the joined table are referenced
// in the query, such as in SELECT DISTINCT or GROUP BY statements.
//
// !!! To ensure dependencies are collected as expected, do not hard-code table names.
// Always build tables using Table or TableAliased. For example:
//
//	// GOOD: The dependency of foo will be collected.
//	foo := Table("foo")
//	b.SELECT(sqlf.F("?.id", foo))
//
//	// BAD: The dependency of foo will NOT be collected.
//	b.SELECT(sqlf.F("foo.id"))
func (b *QueryBuilder) LeftJoinOptional(t TableAliased, on *sqlf.Fragment) *QueryBuilder {
	return b.join("LEFT JOIN", t, on, true)
}

// RightJoin append / replace a right join table.
func (b *QueryBuilder) RightJoin(t TableAliased, on *sqlf.Fragment) *QueryBuilder {
	return b.join("RIGHT JOIN", t, on, false)
}

// FullJoin append / replace a full join table.
func (b *QueryBuilder) FullJoin(t TableAliased, on *sqlf.Fragment) *QueryBuilder {
	return b.join("FULL JOIN", t, on, false)
}

// CrossJoin append / replace a cross join table.
func (b *QueryBuilder) CrossJoin(t TableAliased) *QueryBuilder {
	return b.join("CROSS JOIN", t, nil, false)
}

// join append or replace a join table.
func (b *QueryBuilder) join(joinStr string, t TableAliased, on *sqlf.Fragment, optional bool) *QueryBuilder {
	if t.Name == "" {
		b.pushError(fmt.Errorf("join table name is empty"))
		return b
	}
	// if _, ok := b.tablesDict[t.AppliedName()]; ok {
	// 	if t.Alias == "" {
	// 		b.pushError(fmt.Errorf("table [%s] is already joined", t.Name))
	// 		return b
	// 	}
	// 	b.pushError(fmt.Errorf("table [%s AS %s] is already joined", t.Name, t.Alias))
	// 	return b
	// }
	if len(b.tables) == 0 {
		// reserve the first alias for the main table
		b.tables = append(b.tables, &fromTable{})
	}
	table := &fromTable{
		Names: t,
		Builder: sqlf.F(
			joinStr+" ? ?",
			NewTableAsBuilder(t),
			on.WithPrefix("ON"),
		),
		Optional: optional,
	}
	if target, replacing := b.tablesDict[t.AppliedName()]; replacing {
		*target = *table
		return b
	}
	b.tables = append(b.tables, table)
	b.tablesDict[t.AppliedName()] = table
	return b
}
