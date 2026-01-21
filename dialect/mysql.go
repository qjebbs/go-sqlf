package dialect

var _ Dialect = MySQL{}

// MySQL is the ANSI SQL dialect.
type MySQL struct {
	// BindVar is the bind variable style to use.
	// If empty, StyleQuestion is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleBacktick is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d MySQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d MySQL) QuoteIdentifier(name string) string {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleBacktick.QuoteIdentifier(name)
	}
	return d.Quote.QuoteIdentifier(name)
}
