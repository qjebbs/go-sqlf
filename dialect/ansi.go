package dialect

import (
	"strings"
	"time"
)

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d AnsiSQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d AnsiSQL) QuoteIdentifier(name string) string {
	return IdentifierQuoteStyleDoubleQuote.Quote(name)
}

// FormatTime formats time strings for the dialect.
func (d AnsiSQL) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

// QuoteString quotes a string for use in a query.
func (d AnsiSQL) QuoteString(s string) string {
	var b strings.Builder
	b.WriteString("'")
	for _, r := range s {
		if r == '\x00' {
			// Prevent query truncation attack,
			// repsect it as end of string.
			break
		}
		if r == '\'' {
			b.WriteString("''")
			continue
		}
		b.WriteRune(r)
	}
	b.WriteString("'")
	return b.String()
}
