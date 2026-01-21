package dialect

import (
	"fmt"
	"strings"
	"unicode"
)

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct {
	// BindVar is the bind variable style to use.
	// If empty, StyleQuestion is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleBacktick is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d MySQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteStyle returns the identifier quote style for the dialect.
func (d MySQL) QuoteStyle() QuoteStyle {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleBacktick
	}
	return d.Quote
}

// TimeFormat returns the time format for the dialect.
func (d MySQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}

// QuoteString escapes a string for use in a query, e.g. 'string'.
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
