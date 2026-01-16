package dialect

import (
	"strings"
)

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// Style is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVarStyle BindStyle
}

// BindStyle returns the bind variable style for the dialect.
func (d SQLite) BindStyle() BindStyle {
	if d.BindVarStyle == BindStyleDefault {
		return BindStyleQuestion
	}
	return d.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLite) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// TimeFormat returns the time format for the dialect.
func (d SQLite) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
