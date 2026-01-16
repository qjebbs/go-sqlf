package dialect

import (
	"strings"
)

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct {
	// Style is the bind variable style to use.
	// If empty, StyleQuestion is used.
	BindVarStyle BindStyle
}

// BindStyle returns the bind variable style for the dialect.
func (d MySQL) BindStyle() BindStyle {
	if d.BindVarStyle == BindStyleDefault {
		return BindStyleQuestion
	}
	return d.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d MySQL) QuoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// TimeFormat returns the time format for the dialect.
func (d MySQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
