package arg

import "github.com/qjebbs/go-sqlf/v4/internal/syntax"

// NewArgStoreFromStyle creates a new ArgStore based on the given Style.
func NewArgStoreFromStyle(style syntax.BindStyle) Store {
	switch style {
	case syntax.BindStyleQuestion:
		return NewPositional()
	case syntax.BindStyleDollarNumbered:
		return NewNumbered("$")
	case syntax.BindStyleQuestionNumbered:
		return NewNumbered("?")
	case syntax.BindStyleColonNamed:
		return NewNamed(":", "p")
	case syntax.BindStyleColonNumbered:
		return NewNumbered(":")
	case syntax.BindStyleAtNamed:
		return NewNamed("@", "p")
	default:
		// unknown style, return default
		return NewPositional()
	}
}
