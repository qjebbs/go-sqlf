package dialect

var _ Dialect = SQLServer{}

// SQLServer is the ANSI SQL dialect.
type SQLServer struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleAtNamed is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleSquareBracket is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d SQLServer) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleAtNamed
	}
	return d.BindVar
}

// QuoteIdentifier returns the identifier quote style for the dialect.
func (d SQLServer) QuoteIdentifier(name string) string {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleSquareBracket.QuoteIdentifier(name)
	}
	return d.Quote.QuoteIdentifier(name)
}
