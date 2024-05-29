package sqlb

// TableAliased is the table name with alias.
type TableAliased struct {
	Name, Alias Table
}

// NewTableAliased returns a new TableAliased.
func NewTableAliased(name, alias Table) TableAliased {
	return TableAliased{
		Name:  name,
		Alias: alias,
	}
}

// WithAlias returns a new Table with updated alias.
func (t TableAliased) WithAlias(alias Table) TableAliased {
	return TableAliased{
		Name:  t.Name,
		Alias: alias,
	}
}

// AppliedName returns the alias if it is not empty, otherwise returns the name.
func (t TableAliased) AppliedName() Table {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

// Names returns the table name and alias.
func (t TableAliased) Names() []Table {
	return []Table{t.Name, t.Alias}
}

// Column returns a column of the table.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := NewTable("table", "t")
//	// these two are equivalent
//	t.Column("id")         // "t.id"
//	t.Expression("#f1.id") // "t.id"
//
// If you want to use the column name directly, try:
//
//	t.Expressions("id") // "id"
func (t TableAliased) Column(name string) *Column {
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
//	t.Expressions("#f1.id", "#f1.name")  // "t.id", "t.name"
//
// If you want to use the column name directly, try:
//
//	t.Expressions("id", "name") // "id", "name"
func (t TableAliased) Columns(names ...string) []*Column {
	return t.AppliedName().Columns(names...)
}
