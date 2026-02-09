package dialect

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

var _ Dialect = SQLServer{}

// SQLServer is the ANSI SQL dialect.
type SQLServer struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleAtNamed is used.
	BindVar BindVarStyle
	// IdentifierQuote is the identifier quote style to use.
	// If zero, IdentifierQuoteStyleSquareBracket is used.
	IdentifierQuote IdentifierQuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLServer) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleAtNamed
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d SQLServer) QuoteIdentifier(name string) string {
	if d.IdentifierQuote == IdentifierQuoteStyleDefault {
		return IdentifierQuoteStyleSquareBracket.Quote(name)
	}
	return d.IdentifierQuote.Quote(name)
}

// FormatTime formats time strings for the dialect.
func (d SQLServer) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

// QuoteString quotes a string for use in a query.
func (d SQLServer) QuoteString(s string) string {
	if s == "" {
		return "''"
	}
	var hasUnicode bool
	for _, r := range s {
		if r > 255 {
			hasUnicode = true
			break
		}
	}

	const concatOperator = " + "
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
			if r > 255 {
				b.WriteString("NCHAR(")
			} else {
				b.WriteString("CHAR(")
			}
			b.WriteString(strconv.FormatInt(int64(r), 10))
			b.WriteString(")")
		} else {
			if !inQuote {
				if hasUnicode {
					b.WriteString("N'")
				} else {
					b.WriteString("'")
				}
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
