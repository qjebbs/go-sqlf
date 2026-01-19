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
type TokenType string

const (
	_EOF TokenType = "EOF"

	_Ref     = "ref"
	_Literal = "literal"
	_Escape  = "escape"
	_Plain   = "plain text"
)

type kind uint8

const (
	_KindNone kind = iota
	_KindLitNumber
	_KindLitString

	_KindRefNamed
	_KindRefNumbered
	_KindRefPositional
)
