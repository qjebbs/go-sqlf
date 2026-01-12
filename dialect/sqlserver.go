package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = SQLServer{}

// SQLServer is the ANSI SQL dialect.
type SQLServer struct{}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLServer) QuoteIdentifier(name string) string {
	return `[` + strings.ReplaceAll(name, `]`, `]]`) + `]`
}

// NewArgStore creates a new Positional ArgStore.
func (d SQLServer) NewArgStore() argstore.Store {
	return argstore.NewNamed("@", "p")
}
