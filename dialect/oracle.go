package dialect

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct {
	// BindVar is the bind variable style to use.
	// If empty, StyleColonNumbered is used.
	BindVar BindVarStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d Oracle) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleColonNumbered
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d Oracle) QuoteIdentifier(name string) string {
	return IdentifierQuoteStyleDoubleQuote.Quote(name)
}

// FormatTime formats time strings for the dialect.
func (d Oracle) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

// QuoteString quotes a string for use in a query.
func (d Oracle) QuoteString(s string) string {
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
			b.WriteString("CHR(")
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
