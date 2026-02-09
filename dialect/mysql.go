package dialect

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct {
	// BindVar is the bind variable style to use.
	// If empty, StyleQuestion is used.
	BindVar BindVarStyle
	// IdentifierQuote is the identifier quote style to use.
	// If zero, IdentifierQuoteStyleBacktick is used.
	IdentifierQuote IdentifierQuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d MySQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d MySQL) QuoteIdentifier(name string) string {
	if d.IdentifierQuote == IdentifierQuoteStyleDefault {
		return IdentifierQuoteStyleBacktick.Quote(name)
	}
	return d.IdentifierQuote.Quote(name)
}

// FormatTime formats time strings for the dialect.
func (d MySQL) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

// QuoteString quotes a string for use in a query.
func (d MySQL) QuoteString(s string) string {
	var b strings.Builder
	b.WriteString("'")
	for _, r := range s {
		if r == '\'' {
			b.WriteString("''")
			continue
		}
		isControl := unicode.IsControl(r)
		if !isControl {
			b.WriteRune(r)
		} else {
			switch r {
			case '\b':
				b.WriteString(`\b`)
			case '\f':
				b.WriteString(`\f`)
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			case '\t':
				b.WriteString(`\t`)
			default:
				if r < 256 {
					b.WriteString(fmt.Sprintf("\\x%02x", r))
				} else {
					b.WriteString(fmt.Sprintf("\\u%04x", r))
				}
			}
		}
	}
	b.WriteString("'")
	return b.String()
}
