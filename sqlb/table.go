package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

var _ (sqlf.FragmentBuilder) = Table("")

// Table is a table identifier, it can be a table name or an alias.
type Table string

// BuildFragment implements FragmentBuilder
func (t Table) BuildFragment(_ *sqlf.Context) (query string, err error) {
	return string(t), nil
}

// Column returns a column of the table.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := Table("t")
//	t.Column("id")  // "t.id"
func (t Table) Column(name string) *Column {
	return &Column{
		fragment: sqlf.F(preBuildColumn(t, name)),
		table:    t,
	}
}

func preBuildColumn(t Table, name string) string {
	if name == "" {
		return ""
	}
	if t == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", t, name)
}

// Columns returns columns of the table from names.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := Table("t")
//	t.Columns("id", "name")  // "t.id", "t.name"
func (t Table) Columns(names ...string) []*Column {
	r := make([]*Column, 0, len(names))
	for _, name := range names {
		r = append(r, t.Column(name))
	}
	return r
}

// AnonymousColumn returns a anonymous column of the table.
// For example:
//
//	t := Table("t")
//	t.AnonymousColumn("id")  // "id"
func (t Table) AnonymousColumn(name string) *Column {
	return &Column{
		fragment: sqlf.F(name),
		table:    t,
	}
}

// AnonymousColumns returns anonymous columns of the table from names.
//
// For example:
//
//	t := Table("t")
//	t.Columns("id", "name")  // "id", "name"
func (t Table) AnonymousColumns(names ...string) []*Column {
	r := make([]*Column, 0, len(names))
	for _, name := range names {
		r = append(r, t.AnonymousColumn(name))
	}
	return r
}
