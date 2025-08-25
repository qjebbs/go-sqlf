package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3"
)

var _ (sqlf.Builder) = TableAliased{}

// Build implements sqlf.Builder
func (t TableAliased) Build(ctx *sqlf.Context) (query string, err error) {
	return t.AppliedName().Build(ctx)
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
func (t TableAliased) Column(name string) sqlf.Builder {
	return t.AppliedName().Column(name)
}

// Columns returns columns of the table from names.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := NewTable("table", "t")
//	t.Columns("id", "name")   // "t.id", "t.name"
func (t TableAliased) Columns(names ...string) []sqlf.Builder {
	return t.AppliedName().Columns(names...)
}
