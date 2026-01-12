package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLite) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d SQLite) NewArgStore() argstore.Store {
	return argstore.NewPositional()
}
