package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/arg"
)

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct {
	// Style is the bind variable style to use.
	// If zero, StyleDollarNumbered is used.
	BindVarStyle arg.Style
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d PostgreSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d PostgreSQL) NewArgStore() arg.Store {
	return arg.NewArgStoreFromStyle(d.BindVarStyle, arg.StyleDollarNumbered)
}

// TimeFormat returns the time format for the dialect.
func (d PostgreSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999-07:00"
}
