package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3"
)

var _ (sqlf.Builder) = Table("")

// Table is a table identifier, it can be a table name or an alias.
type Table string

// Build implements sqlf.Builder
func (t Table) Build(ctx *sqlf.Context) (query string, err error) {
	if v := ctx.Value(depTablesKey{}); v != nil {
		if deps, ok := v.(map[Table]bool); ok && deps != nil {
			// collecting
			deps[t] = true
		}
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
func (t Table) Column(name string) sqlf.Builder {
	return sqlf.F("?."+name, t)
}

// Columns returns columns of the table from names.
// It adds table prefix to the column name, e.g.: "id" -> "t.id".
//
// For example:
//
//	t := Table("t")
//	t.Columns("id", "name")  // "t.id", "t.name"
func (t Table) Columns(names ...string) []sqlf.Builder {
	r := make([]sqlf.Builder, 0, len(names))
	for _, name := range names {
		r = append(r, t.Column(name))
	}
	return r
}
