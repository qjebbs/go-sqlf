package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/argstore"
)

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// BindVarStyle is the bind variable style to use.
	// If empty, "?" is used.
	BindVarStyle SQLiteBindVarStyle
}

// SQLiteBindVarStyle is the bind variable style to use.
type SQLiteBindVarStyle int

const (
	// SQLiteBindVarStyleDefault is the default bind variable style.
	SQLiteBindVarStyleDefault SQLiteBindVarStyle = iota
	// SQLiteBindVarStyleQuestion is the "?" bind variable style.
	SQLiteBindVarStyleQuestion
	// SQLiteBindVarStyleQuestionNumbered is the "?NNN" bind variable style.
	SQLiteBindVarStyleQuestionNumbered
	// SQLiteBindVarStyleColonNamed is the ":AAAA" bind variable style.
	SQLiteBindVarStyleColonNamed
	// SQLiteBindVarStyleAtNamed is the "@AAAA" bind variable style.
	SQLiteBindVarStyleAtNamed
	// SQLiteBindVarStyleDollarNumbered is the "$AAAA" bind variable style.
	SQLiteBindVarStyleDollarNumbered
)

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLite) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// NewArgStore creates a new Positional ArgStore.
func (d SQLite) NewArgStore() argstore.Store {
	switch d.BindVarStyle {
	case SQLiteBindVarStyleQuestionNumbered:
		return argstore.NewNumbered("?")
	case SQLiteBindVarStyleColonNamed:
		return argstore.NewNamed(":", "")
	case SQLiteBindVarStyleAtNamed:
		return argstore.NewNamed("@", "")
	case SQLiteBindVarStyleDollarNumbered:
		return argstore.NewNumbered("$")
	case SQLiteBindVarStyleQuestion:
		fallthrough
	case SQLiteBindVarStyleDefault:
		fallthrough
	default:
		return argstore.NewPositional()
	}
}

// TimeFormat returns the time format for the dialect.
func (d SQLite) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
