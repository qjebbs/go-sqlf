package argstore

// BindVarStyle is the bind variable style to use.
type BindVarStyle int

const (
	// BindVarStyleDefault is the default bind variable style.
	BindVarStyleDefault BindVarStyle = iota
	// BindVarStyleQuestion is the "?" bind variable style.
	BindVarStyleQuestion
	// BindVarStyleDollarNumbered is the "$1" bind variable style.
	BindVarStyleDollarNumbered
	// BindVarStyleQuestionNumbered is the "?1" bind variable style.
	BindVarStyleQuestionNumbered
	// BindVarStyleColonNamed is the ":name" bind variable style.
	BindVarStyleColonNamed
	// BindVarStyleColonNumbered is the ":1" bind variable style.
	BindVarStyleColonNumbered
	// BindVarStyleAtNamed is the "@name" bind variable style.
	BindVarStyleAtNamed
)

// NewArgStoreFromBindVarStyle creates a new ArgStore based on the given BindVarStyle.
func NewArgStoreFromBindVarStyle(style, defaultStyle BindVarStyle) Store {
	if style == BindVarStyleDefault {
		style = defaultStyle
	}
	switch style {
	case BindVarStyleQuestion:
		return NewPositional()
	case BindVarStyleDollarNumbered:
		return NewNumbered("$")
	case BindVarStyleQuestionNumbered:
		return NewNumbered("?")
	case BindVarStyleColonNamed:
		return NewNamed(":", "p")
	case BindVarStyleColonNumbered:
		return NewNumbered(":")
	case BindVarStyleAtNamed:
		return NewNamed("@", "p")
	default:
		// unknown style, return default
		return NewPositional()
	}
}
