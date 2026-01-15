package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/arg"
)

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct {
	// Style is the bind variable style to use.
	// If empty, StyleColonNumbered is used.
	BindVarStyle arg.Style
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d Oracle) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d Oracle) NewArgStore() arg.Store {
	return arg.NewArgStoreFromStyle(d.BindVarStyle, arg.StyleColonNumbered)
}

// TimeFormat returns the time format for the dialect.
func (d Oracle) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
