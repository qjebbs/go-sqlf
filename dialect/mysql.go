package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/arg"
)

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct {
	// Style is the bind variable style to use.
	// If empty, StyleQuestion is used.
	BindVarStyle arg.Style
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d MySQL) QuoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// NewArgStore creates a new Positional ArgStore.
func (d MySQL) NewArgStore() arg.Store {
	return arg.NewArgStoreFromStyle(d.BindVarStyle, arg.StyleQuestion)
}

// TimeFormat returns the time format for the dialect.
func (d MySQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
