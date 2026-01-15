package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/arg"
)

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// Style is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVarStyle arg.Style
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLite) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d SQLite) NewArgStore() arg.Store {
	return arg.NewArgStoreFromStyle(d.BindVarStyle, arg.StyleQuestion)
}

// TimeFormat returns the time format for the dialect.
func (d SQLite) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
