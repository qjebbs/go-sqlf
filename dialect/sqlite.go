package dialect

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLite) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteStyle returns the identifier quote style for the dialect.
func (d SQLite) QuoteStyle() QuoteStyle {
	return QuoteStyleDoubleQuote
}

// TimeFormat returns the time format for the dialect.
func (d SQLite) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
