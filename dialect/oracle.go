package dialect

var _ Dialect = Oracle{}

// Oracle is the ANSI SQL dialect.
type Oracle struct {
	// BindVar is the bind variable style to use.
	// If empty, StyleColonNumbered is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleDoubleQuote is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d Oracle) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleColonNumbered
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d Oracle) QuoteIdentifier(name string) string {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleDoubleQuote.QuoteIdentifier(name)
	}
	return d.Quote.QuoteIdentifier(name)
}
