package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

// From set the from table.
func (b *QueryBuilder) From(t TableAliased) *QueryBuilder {
	if t.Name == "" {
		b.pushError(fmt.Errorf("from table is empty"))
		return b
	}
	tableAndAlias := string(t.Name)
	if t.Alias != "" {
		tableAndAlias = tableAndAlias + " AS " + string(t.Alias)
	}
	table := &fromTable{
		Names:    t,
		Fragment: sqlf.F(tableAndAlias),
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

// LeftJoinOptional append / replace a left join table, and mark it as optional.
//
// CAUSION:
//
//   - Make sure all columns referenced in the query are reflected in
//     *sqlf.Fragment.Columns, so that the *QueryBuilder can calculate the dependency
//     correctly.
//   - Make sure it's used with the SELECT DISTINCT statement, otherwise it works
//     exactly the same as LeftJoin().
//
// Consider the following two queries:
//
//	SELECT DISTINCT foo.* FROM foo LEFT JOIN bar ON foo.id = bar.foo_id
//	SELECT DISTINCT foo.* FROM foo
//
// They return the same result, but the second query more efficient.
// If the join to "bar" is declared with LeftJoinOptional(), *QueryBuilder
// will trim it if no relative columns referenced in the query, aka Join Elimination.
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
	tableAndAlias := t.Name
	if t.Alias != "" {
		tableAndAlias = tableAndAlias + " AS " + t.Alias
	}
	table := &fromTable{
		Names: t,
		Fragment: sqlf.Ff(
			fmt.Sprintf("%s %s #f1", joinStr, tableAndAlias),
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
