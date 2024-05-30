package sqlb

import "github.com/qjebbs/go-sqlf/v2"

var _ (sqlf.FragmentBuilder) = TableAliased{}

// BuildFragment implements FragmentBuilder
func (t TableAliased) BuildFragment(_ *sqlf.Context) (query string, err error) {
	return string(t.AppliedName()), nil
}

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
//	t.Column("id")  // "t.id"
func (t TableAliased) Column(name string) *Column {
	return t.AppliedName().Column(name)
}

// Columns returns columns of the table from names.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := NewTable("table", "t")
//	t.Columns("id", "name")   // "t.id", "t.name"
func (t TableAliased) Columns(names ...string) []*Column {
	return t.AppliedName().Columns(names...)
}

// AnonymousColumn returns a anonymous column of the table.
// For example:
//
//	t := NewTable("table", "t")
//	t.AnonymousColumn("id")  // "id"
func (t TableAliased) AnonymousColumn(name string) *Column {
	return t.AppliedName().AnonymousColumn(name)
}

// AnonymousColumns returns anonymous columns of the table from names.
//
// For example:
//
//	t := NewTable("table", "t")
//	t.Columns("id", "name")  // "id", "name"
func (t TableAliased) AnonymousColumns(names ...string) []*Column {
	return t.AppliedName().AnonymousColumns(names...)
}
