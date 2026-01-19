package arg

import "github.com/qjebbs/go-sqlf/v4/internal/syntax"

// NewArgStoreFromStyle creates a new ArgStore based on the given Style.
func NewArgStoreFromStyle(style syntax.BindVarStyle) Store {
	switch style {
	case syntax.BindVarStyleQuestion:
		return NewPositional()
	case syntax.BindVarStyleDollarNumbered:
		return NewNumbered("$")
	case syntax.BindVarStyleQuestionNumbered:
		return NewNumbered("?")
	case syntax.BindVarStyleColonNamed:
		return NewNamed(":", "p")
	case syntax.BindVarStyleColonNumbered:
		return NewNumbered(":")
	case syntax.BindVarStyleAtNamed:
		return NewNamed("@", "p")
	default:
		// unknown style, return default
		return NewPositional()
	}
}
