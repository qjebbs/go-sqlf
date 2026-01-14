package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct {
	// BindVarStyle is the bind variable style to use.
	// If zero, BindVarStyleDollarNumbered is used.
	BindVarStyle argstore.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d PostgreSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d PostgreSQL) NewArgStore() argstore.Store {
	return argstore.NewArgStoreFromBindVarStyle(d.BindVarStyle, argstore.BindVarStyleDollarNumbered)
}

// TimeFormat returns the time format for the dialect.
func (d PostgreSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999-07:00"
}
