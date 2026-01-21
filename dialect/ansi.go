package dialect

import "strings"

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleDoubleQuote is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d AnsiSQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteStyle returns the identifier quote style for the dialect.
func (d AnsiSQL) QuoteStyle() QuoteStyle {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleDoubleQuote
	}
	return d.Quote
}

// TimeFormat returns the time format for the dialect.
func (d AnsiSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}

// QuoteString escapes a string for use in a query, e.g. 'string'.
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
