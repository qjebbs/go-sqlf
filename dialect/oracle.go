package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d Oracle) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d Oracle) NewArgStore() argstore.Store {
	return argstore.NewNumbered(":")
}
