package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct {
	// BindVarStyle is the bind variable style to use.
	// If empty, ":1" is used.
	BindVarStyle OracleBindVarStyle
}

// OracleBindVarStyle is the bind variable style to use.
type OracleBindVarStyle int

const (
	// OracleBindVarStyleDefault is the default bind variable style.
	OracleBindVarStyleDefault OracleBindVarStyle = iota
	// OracleBindVarStyleNumbered is the ":1" bind variable style.
	OracleBindVarStyleNumbered
	// OracleBindVarStyleNamed is the ":name" bind variable style.
	OracleBindVarStyleNamed
)

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d Oracle) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d Oracle) NewArgStore() argstore.Store {
	switch d.BindVarStyle {
	case OracleBindVarStyleNamed:
		return argstore.NewNamed(":", "p")
	case OracleBindVarStyleNumbered:
		return argstore.NewNumbered(":")
	default:
		return argstore.NewNumbered(":")
	}
}

// TimeFormat returns the time format for the dialect.
func (d Oracle) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
