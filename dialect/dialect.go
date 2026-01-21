package dialect

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

// Dialect defines SQL dialects.
type Dialect interface {
	// BindVarStyle returns the bind variable style for the dialect.
	BindVarStyle() BindVarStyle
	// QuoteStyle returns the identifier quote style for the dialect.
	QuoteStyle() QuoteStyle
	// TimeFormat returns the time format for the dialect.
	TimeFormat() string
	// QuoteString quotes a string for use in a query.
	//
	// It should escape any special characters as needed by the dialect,
	// and take care of any necessary prefixing (e.g. E'...' in PostgreSQL, N'...' in SQL Server).
	//
	// Examples:
	//   QuoteString("str") // 'str'
	//   SQLite.QuoteString("str\nstr") // 'str' || CHAR(10) || 'str'
	//   PostgreSQL.QuoteString("str\nstr") // E'str\nstr'
	QuoteString(s string) string
}

// BindVarStyle is the bind variable style to use.
type BindVarStyle = syntax.BindVarStyle

const (
	// BindVarStyleDefault is the default bind variable style.
	BindVarStyleDefault BindVarStyle = syntax.BindVarStyleUnknown
	// BindVarStyleQuestion is the "?" bind variable style.
	BindVarStyleQuestion BindVarStyle = syntax.BindVarStyleQuestion
	// BindVarStyleDollarNumbered is the "$1" bind variable style.
	BindVarStyleDollarNumbered BindVarStyle = syntax.BindVarStyleDollarNumbered
	// BindVarStyleQuestionNumbered is the "?1" bind variable style.
	BindVarStyleQuestionNumbered BindVarStyle = syntax.BindVarStyleQuestionNumbered
	// BindVarStyleColonNamed is the ":name" bind variable style.
	BindVarStyleColonNamed BindVarStyle = syntax.BindVarStyleColonNamed
	// BindVarStyleColonNumbered is the ":1" bind variable style.
	BindVarStyleColonNumbered BindVarStyle = syntax.BindVarStyleColonNumbered
	// BindVarStyleAtNamed is the "@name" bind variable style.
	BindVarStyleAtNamed BindVarStyle = syntax.BindVarStyleAtNamed
)

// QuoteStyle is the identifier quote style to use.
type QuoteStyle int

const (
	// QuoteStyleDefault quotes identifiers with the default style.
	QuoteStyleDefault QuoteStyle = iota
	// QuoteStyleDoubleQuote quotes identifiers with double quotes, e.g. "identifier".
	QuoteStyleDoubleQuote
	// QuoteStyleBacktick quotes identifiers with backticks, e.g. `identifier`.
	QuoteStyleBacktick
	// QuoteStyleSquareBracket quotes identifiers with square brackets, e.g. [identifier].
	QuoteStyleSquareBracket
)

// QuoteIdentifier quotes an identifier using the given quote style.
func (s QuoteStyle) QuoteIdentifier(name string) string {
	switch s {
	case QuoteStyleDoubleQuote:
		return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
	case QuoteStyleBacktick:
		return "`" + strings.ReplaceAll(name, "`", "``") + "`"
	case QuoteStyleSquareBracket:
		return `[` + strings.ReplaceAll(name, `]`, `]]`) + `]`
	default:
		return name
	}
}
