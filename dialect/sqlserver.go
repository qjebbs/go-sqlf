package dialect

import (
	"strconv"
	"strings"
	"unicode"
)

var _ Dialect = SQLServer{}

// SQLServer is the ANSI SQL dialect.
type SQLServer struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleAtNamed is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleSquareBracket is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLServer) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleAtNamed
	}
	return d.BindVar
}

// QuoteStyle returns the identifier quote style for the dialect.
func (d SQLServer) QuoteStyle() QuoteStyle {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleSquareBracket
	}
	return d.Quote
}

// TimeFormat returns the time format for the dialect.
func (d SQLServer) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}

// QuoteString escapes a string for use in a query, e.g. 'string'.
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
