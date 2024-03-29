package sqlb

import "github.com/qjebbs/go-sqlf"

// Table is the table name with alias.
type Table struct {
	Name, Alias sqlf.Table
}

// NewTable returns a new Table.
func NewTable(name, alias sqlf.Table) Table {
	return Table{
		Name:  name,
		Alias: alias,
	}
}

// WithAlias returns a new Table with updated alias.
func (t Table) WithAlias(alias sqlf.Table) Table {
	return Table{
		Name:  t.Name,
		Alias: alias,
	}
}

// AppliedName returns the alias if it is not empty, otherwise returns the name.
func (t Table) AppliedName() sqlf.Table {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

// Names returns the table name and alias.
func (t Table) Names() []sqlf.Table {
	return []sqlf.Table{t.Name, t.Alias}
}

// Column returns a column of the table.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := NewTable("table", "t")
//	// these two are equivalent
//	t.Column("id")         // "t.id"
//	t.Expression("#t1.id") // "t.id"
//
// If you want to use the column name directly, try:
//
//	t.Expressions("id") // "id"
func (t Table) Column(name string) *sqlf.Column {
	return t.AppliedName().Column(name)
}

// Columns returns columns of the table from names.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := NewTable("table", "t")
//	// these two are equivalent
//	t.Columns("id", "name")              // "t.id", "t.name"
//	t.Expressions("#t1.id", "#t1.name")  // "t.id", "t.name"
//
// If you want to use the column name directly, try:
//
//	t.Expressions("id", "name") // "id", "name"
func (t Table) Columns(names ...string) []*sqlf.Column {
	return t.AppliedName().Columns(names...)
}

// Expression returns a column of the table from the expression, it accepts
// bindvars and the preprocessor #t1 which is implicit in t.AppliedName().
//
// For example:
//
//	t := NewTable("table", "t")
//	t.Expression("id")                      // "id"
//	t.Expression("#t1.id")                  // "table.id"
//	t.Expression("#t1.id")                  // "t.id"
//	t.Expression("COALESCE(#t1.id,0)")      // "COALESCE(t.id,0)"
//	t.Expression("#t1.deteled_at > $1", 1)  // "t.deteled_at > $1"
func (t Table) Expression(expression string, args ...any) *sqlf.Column {
	return t.AppliedName().Expression(expression, args...)
}

// Expressions returns columns of the table from the expression, it accepts
// bindvars and the preprocessor #t1 which is implicit in t.AppliedName().
//
// For example:
//
//	t := NewTable("table", "t")
//	t.Expressions("id", "deteled_at")         // "id", "deteled_at"
//	t.Expressions("#t1.id", "#t1.deteled_at") // "table.id", "table.deteled_at"
//	t.Expressions("#t1.id", "#t1.deteled_at") // "t.id", "t.deteled_at"
//	t.Expressions("COALESCE(#t1.id,0)")       // "COALESCE(t.id,0)"
func (t Table) Expressions(expressions ...string) []*sqlf.Column {
	return t.AppliedName().Expressions(expressions...)
}
