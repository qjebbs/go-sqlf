package dialect

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleDollarNumbered is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleDoubleQuote is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d PostgreSQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleDollarNumbered
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d PostgreSQL) QuoteIdentifier(name string) string {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleDoubleQuote.QuoteIdentifier(name)
	}
	return d.Quote.QuoteIdentifier(name)
}
