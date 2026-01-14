package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct {
	// BindVarStyle is the bind variable style to use.
	// If zero, BindVarStyleQuestion is used.
	BindVarStyle argstore.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d AnsiSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d AnsiSQL) NewArgStore() argstore.Store {
	return argstore.NewArgStoreFromBindVarStyle(d.BindVarStyle, argstore.BindVarStyleQuestion)
}

// TimeFormat returns the time format for the dialect.
func (d AnsiSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
