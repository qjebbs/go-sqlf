package arg

// Style is the bind variable style to use.
type Style int

const (
	// StyleDefault is the default bind variable style.
	StyleDefault Style = iota
	// StyleQuestion is the "?" bind variable style.
	StyleQuestion
	// StyleDollarNumbered is the "$1" bind variable style.
	StyleDollarNumbered
	// StyleQuestionNumbered is the "?1" bind variable style.
	StyleQuestionNumbered
	// StyleColonNamed is the ":name" bind variable style.
	StyleColonNamed
	// StyleColonNumbered is the ":1" bind variable style.
	StyleColonNumbered
	// StyleAtNamed is the "@name" bind variable style.
	StyleAtNamed
)

// NewArgStoreFromStyle creates a new ArgStore based on the given Style.
func NewArgStoreFromStyle(style, defaultStyle Style) Store {
	if style == StyleDefault {
		style = defaultStyle
	}
	switch style {
	case StyleQuestion:
		return NewPositional()
	case StyleDollarNumbered:
		return NewNumbered("$")
	case StyleQuestionNumbered:
		return NewNumbered("?")
	case StyleColonNamed:
		return NewNamed(":", "p")
	case StyleColonNumbered:
		return NewNumbered(":")
	case StyleAtNamed:
		return NewNamed("@", "p")
	default:
		// unknown style, return default
		return NewPositional()
	}
}
