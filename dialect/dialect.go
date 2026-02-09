package dialect

import (
	"strings"
	"time"

	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

// Dialect defines SQL dialects.
type Dialect interface {
	// BindVarStyle returns the bind variable style for the dialect.
	BindVarStyle() BindVarStyle
	// QuoteIdentifier returns the quoted identifier for the dialect.
	// Examples:
	//   PostgreSQL.QuoteIdentifier("identifier") // "identifier"
	//   MySQL.QuoteIdentifier("identifier") // `identifier`
	//   SQLServer.QuoteIdentifier("identifier") // [identifier]
	QuoteIdentifier(name string) string
	// QuoteString quotes a string for use in a query.
	// It's used in Interpolate() only.
	//
	// It should escape any special characters as needed by the dialect,
	// and take care of any necessary prefixing (e.g. E'...' in PostgreSQL, N'...' in SQL Server).
	//
	// Examples:
	//   QuoteString("str") // 'str'
	//   SQLite.QuoteString("str\nstr") // 'str' || CHAR(10) || 'str'
	//   PostgreSQL.QuoteString("str\nstr") // E'str\nstr'
	QuoteString(s string) string
	// FormatTime formats time strings for the dialect.
	FormatTime(t time.Time) string
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

// IdentifierQuoteStyle is the identifier quote style to use.
type IdentifierQuoteStyle int

const (
	// IdentifierQuoteStyleDefault quotes identifiers with the default style.
	IdentifierQuoteStyleDefault IdentifierQuoteStyle = iota
	// IdentifierQuoteStyleDoubleQuote quotes identifiers with double quotes, e.g. "identifier".
	IdentifierQuoteStyleDoubleQuote
	// IdentifierQuoteStyleBacktick quotes identifiers with backticks, e.g. `identifier`.
	IdentifierQuoteStyleBacktick
	// IdentifierQuoteStyleSquareBracket quotes identifiers with square brackets, e.g. [identifier].
	IdentifierQuoteStyleSquareBracket
)

// Quote quotes an identifier using the given quote style.
func (s IdentifierQuoteStyle) Quote(name string) string {
	switch s {
	case IdentifierQuoteStyleDoubleQuote:
		return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
	case IdentifierQuoteStyleBacktick:
		return "`" + strings.ReplaceAll(name, "`", "``") + "`"
	case IdentifierQuoteStyleSquareBracket:
		return `[` + strings.ReplaceAll(name, `]`, `]]`) + `]`
	default:
		return name
	}
}
