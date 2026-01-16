package dialect

import (
	"strings"
)

var _ Dialect = SQLServer{}

// SQLServer is the ANSI SQL dialect.
type SQLServer struct {
	// Style is the bind variable style to use.
	// If zero, StyleAtNamed is used.
	BindVarStyle BindStyle
}

// BindStyle returns the bind variable style for the dialect.
func (d SQLServer) BindStyle() BindStyle {
	if d.BindVarStyle == BindStyleDefault {
		return BindStyleAtNamed
	}
	return d.BindVarStyle
}

// QuoteIdentifier quotes an identifier using ANSI SQL standard.
func (d SQLServer) QuoteIdentifier(name string) string {
	return `[` + strings.ReplaceAll(name, `]`, `]]`) + `]`
}

// TimeFormat returns the time format for the dialect.
func (d SQLServer) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
