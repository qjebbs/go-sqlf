package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d PostgreSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d PostgreSQL) NewArgStore() argstore.Store {
	return argstore.NewNumbered("$")
}
