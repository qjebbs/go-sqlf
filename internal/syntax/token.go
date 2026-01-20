package syntax

type token struct {
	typ   TokenType
	lit   string
	bad   bool
	kind  kind
	start int
	end   int
	pos   Pos
}

// TokenType is the type of token.
type TokenType uint8

const (
	_EOF TokenType = iota

	_BindVar
	_Literal
	_Escape
	_Raw
)

var tokenTypeStrings = [...]string{
	_EOF: "EOF",

	_BindVar: "BindVar",
	_Literal: "Literal",
	_Escape:  "Escape",
	_Raw:     "Raw",

	"Unknown",
}

func (t TokenType) String() string {
	if t > _Raw {
		return tokenTypeStrings[len(tokenTypeStrings)-1]
	}
	return tokenTypeStrings[t]
}

type kind uint8

const (
	_KindNone kind = iota

	_KindLitNumber
	_KindLitString

	_KindBindVarNamed
	_KindBindVarNumbered
	_KindBindVarPositional
)

var kindStrings = [...]string{
	_KindNone: "None",

	_KindLitNumber: "LitNumber",
	_KindLitString: "LitString",

	_KindBindVarNamed:      "BindVarNamed",
	_KindBindVarNumbered:   "BindVarNumbered",
	_KindBindVarPositional: "BindVarPositional",

	"Unknown",
}

func (k kind) String() string {
	if k > _KindBindVarPositional {
		return kindStrings[len(kindStrings)-1]
	}
	return kindStrings[k]
}
