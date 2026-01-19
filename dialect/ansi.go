package dialect

var _ Dialect = AnsiSQL{}

// AnsiSQL is the ANSI SQL dialect.
type AnsiSQL struct {
	// BindVar is the bind variable style to use.
	// If zero, StyleQuestion is used.
	BindVar BindVarStyle
	// Quote is the identifier quote style to use.
	// If zero, StyleDoubleQuote is used.
	Quote QuoteStyle
}

// BindVarStyle returns the bind variable style for the dialect.
func (d AnsiSQL) BindVarStyle() BindVarStyle {
	if d.BindVar == BindVarStyleDefault {
		return BindVarStyleQuestion
	}
	return d.BindVar
}

// QuoteStyle returns the identifier quote style for the dialect.
func (d AnsiSQL) QuoteStyle() QuoteStyle {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleDoubleQuote
	}
	return d.Quote
}

// TimeFormat returns the time format for the dialect.
func (d AnsiSQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
