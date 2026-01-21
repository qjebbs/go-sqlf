package dialect

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d AnsiSQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d AnsiSQL) QuoteIdentifier(name string) string {
	return IdentifierQuoteStyleDoubleQuote.Quote(name)
}
