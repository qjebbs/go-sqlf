package dialect

var _ Dialect = SQLite{}

// SQLite is the ANSI SQL dialect.
type SQLite struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
	// IdentifierQuote is the identifier quote style to use.
	// If zero, IdentifierQuoteStyleDoubleQuote is used.
	IdentifierQuote IdentifierQuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLite) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d SQLite) QuoteIdentifier(name string) string {
	if d.IdentifierQuote == IdentifierQuoteStyleDefault {
		return IdentifierQuoteStyleDoubleQuote.Quote(name)
	}
	return d.IdentifierQuote.Quote(name)
}
