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

// QuoteStyle returns the identifier quote style for the dialect.
func (d MySQL) QuoteStyle() QuoteStyle {
	if d.Quote == QuoteStyleDefault {
		return QuoteStyleBacktick
	}
	return d.Quote
}

// TimeFormat returns the time format for the dialect.
func (d MySQL) TimeFormat() string {
	return "2006-01-02 15:04:05.999"
}
