package dialect

var _ Dialect = PostgreSQL{}

// PostgreSQL is the ANSI SQL dialect.
type PostgreSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleDollarNumbered is used.
	BindVar BindVarStyle
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
	return IdentifierQuoteStyleDoubleQuote.Quote(name)
}
