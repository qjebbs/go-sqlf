package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d AnsiSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d AnsiSQL) NewArgStore() argstore.Store {
	return argstore.NewPositional()
}
