package dialect

import (
	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

// Dialect defines SQL dialects.
type Dialect interface {
	// BindStyle returns the bind variable style for the dialect.
	BindStyle() BindStyle
	// QuoteIdentifier quotes an identifier based on the dialect.
	QuoteIdentifier(name string) string
	// TimeFormat returns the time format for the dialect.
	TimeFormat() string
}

// BindStyle is the bind variable style to use.
type BindStyle = syntax.BindStyle

const (
	// BindStyleDefault is the default bind variable style.
	BindStyleDefault BindStyle = syntax.BindStyleUnknown
	// BindStyleQuestion is the "?" bind variable style.
	BindStyleQuestion BindStyle = syntax.BindStyleQuestion
	// BindStyleDollarNumbered is the "$1" bind variable style.
	BindStyleDollarNumbered BindStyle = syntax.BindStyleDollarNumbered
	// BindStyleQuestionNumbered is the "?1" bind variable style.
	BindStyleQuestionNumbered BindStyle = syntax.BindStyleQuestionNumbered
	// BindStyleColonNamed is the ":name" bind variable style.
	BindStyleColonNamed BindStyle = syntax.BindStyleColonNamed
	// BindStyleColonNumbered is the ":1" bind variable style.
	BindStyleColonNumbered BindStyle = syntax.BindStyleColonNumbered
	// BindStyleAtNamed is the "@name" bind variable style.
	BindStyleAtNamed BindStyle = syntax.BindStyleAtNamed
)
