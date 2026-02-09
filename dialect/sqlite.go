package dialect

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
	// IdentifierQuote is the identifier quote style to use.
	// If zero, IdentifierQuoteStyleDoubleQuote is used.
	IdentifierQuote IdentifierQuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLite) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d SQLite) QuoteIdentifier(name string) string {
	if d.IdentifierQuote == IdentifierQuoteStyleDefault {
		return IdentifierQuoteStyleDoubleQuote.Quote(name)
	}
	return d.IdentifierQuote.Quote(name)
}

// FormatTime formats time strings for the dialect.
func (d SQLite) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

// QuoteString quotes a string for use in a query.
func (d SQLite) QuoteString(s string) string {
	if s == "" {
		return "''"
	}
	const concatOperator = " || "
	var b strings.Builder
	inQuote := false
	for i, r := range s {
		if i > 0 && !inQuote {
			b.WriteString(concatOperator)
		}
		isControl := unicode.IsControl(r)
		if isControl {
			if inQuote {
				b.WriteString("'")
				b.WriteString(concatOperator)
				inQuote = false
			}
			b.WriteString("CHAR(")
			b.WriteString(strconv.FormatInt(int64(r), 10))
			b.WriteString(")")
		} else {
			if !inQuote {
				b.WriteString("'")
				inQuote = true
			}
			if r == '\'' {
				b.WriteString("''")
			} else {
				b.WriteRune(r)
			}
		}
	}
	if inQuote {
		b.WriteString("'")
	}
	return b.String()
}
