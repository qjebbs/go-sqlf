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

func (t TokenType) String() string {
	switch t {
	case _EOF:
		return "EOF"
	case _BindVar:
		return "BindVar"
	case _Literal:
		return "Literal"
	case _Escape:
		return "Escape"
	case _Raw:
		return "Raw"
	}
	return "Unknown"
}

type kind uint8

const (
	_KindNone kind = iota
	_KindLitNumber
	_KindLitString

	_KindRefNamed
	_KindRefNumbered
	_KindRefPositional
)
