package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d MySQL) QuoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// NewArgStore creates a new Positional ArgStore.
func (d MySQL) NewArgStore() argstore.Store {
	return argstore.NewPositional()
}

// TimeFormat returns the time format for the dialect.
func (d MySQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
