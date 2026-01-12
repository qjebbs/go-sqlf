package dialect

import (
	"github.com/qjebbs/go-sqlf/v4/argstore"
)

// Dialect defines SQL dialects.
type Dialect interface {
	// QuoteIdentifier quotes an identifier based on the dialect.
	QuoteIdentifier(name string) string
	// NewArgStore creates a new ArgStore based on the dialect.
	NewArgStore() argstore.Store
}
