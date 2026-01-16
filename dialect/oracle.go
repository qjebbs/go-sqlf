package dialect

import (
	"strings"
)

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct {
	// Style is the bind variable style to use.
	// If empty, StyleColonNumbered is used.
	BindVarStyle BindStyle
}

// BindStyle returns the bind variable style for the dialect.
func (d Oracle) BindStyle() BindStyle {
	if d.BindVarStyle == BindStyleDefault {
		return BindStyleColonNumbered
	}
	return d.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d Oracle) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// TimeFormat returns the time format for the dialect.
func (d Oracle) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
