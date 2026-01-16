package dialect

import (
	"strings"
)

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct {
	// Style is the bind variable style to use.
	// If zero, StyleDollarNumbered is used.
	BindVarStyle BindStyle
}

// BindStyle returns the bind variable style for the dialect.
func (d PostgreSQL) BindStyle() BindStyle {
	if d.BindVarStyle == BindStyleDefault {
		return BindStyleDollarNumbered
	}
	return d.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d PostgreSQL) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// TimeFormat returns the time format for the dialect.
func (d PostgreSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999-07:00"
}
