package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3"
)

var _ (sqlf.Builder) = Table("")

// Table is a table identifier, it can be a table name or an alias.
type Table string

// Build implements sqlf.Builder
func (t Table) Build(ctx *sqlf.Context) (query string, err error) {
	deps := depsFromContext(ctx)
	if deps != nil {
		// collecting
		deps[t] = true
	}
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
		fragment: sqlf.F("?."+name, t),
	}
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
		fragment:       sqlf.F(name),
		anonymousTable: t,
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
