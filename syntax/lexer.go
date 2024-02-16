package syntax

import (
	"unicode"
	"unicode/utf8"
)

// EOF is the EOF rune
const EOF rune = 0

// lexerHelper is the lexical analyzer
type lexerHelper struct {
	input string // input string

	start   pos  // position where current token starts
	current pos  // current postion of the lexer
	rune    rune // current rune
	width   int  // bytes width of the current rune
}

type pos struct {
	offset int
	Pos
}

// newLexerHelper returns a new LexerHelper
func newLexerHelper(input string) *lexerHelper {
	l := &lexerHelper{
		input:   input,
		current: pos{Pos: NewPos(1, 0)},
	}
	l.Next()
	return l
}

// StartToken was called when emit a new token, who set the l.Start
// for the next token.
func (l *lexerHelper) StartToken(tokens ...any) {
	l.start = l.current
}

// Next moves to the next rune
func (l *lexerHelper) Next() rune {
	if l.current.offset >= len(l.input) {
		l.width = 0
		l.rune = EOF
		return l.rune
	}
	l.current.offset += l.width
	result, width := utf8.DecodeRuneInString(l.input[l.current.offset:])
	l.width = width
	l.rune = result
	if l.rune == '\n' {
		l.current.line++
		l.current.col = 0
	} else {
		l.current.col++
	}
	return l.rune
}

// Peek returns the next rune without changing the postions.
// it returns the whole string left if n is 0.
func (l *lexerHelper) Peek() rune {
	next := l.current.offset + l.width
	if next >= len(l.input) {
		return EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[next:])
	return r
}

// SkipWhitespace skips all leading whitespaces
func (l *lexerHelper) SkipWhitespace() {
	for {
		if !unicode.IsSpace(l.rune) {
			break
		}
		l.Next()
		if l.rune == EOF {
			break
		}
	}
}

// IsEOF tells if it reaches the EOF
func (l *lexerHelper) IsEOF() bool {
	return l.current.offset >= len(l.input)
}

// IsWhitespace tells if it's currently a whitespace
func (l *lexerHelper) IsWhitespace() bool {
	return unicode.IsSpace(l.rune)
}

// Advanced tells if current position is advanced compared to the start position
func (l *lexerHelper) Advanced() bool {
	return l.current.offset > l.start.offset
}

// Lower returns lower-case ch if ch is ASCII letter
func (l *lexerHelper) Lower() rune {
	return ('a' - 'A') | l.rune
}

func (l *lexerHelper) IsLetter() bool {
	return 'a' <= l.Lower() && l.Lower() <= 'z' || l.rune == '_'
}

func (l *lexerHelper) IsDecimal() bool {
	return '0' <= l.rune && l.rune <= '9'
}

func (l *lexerHelper) IsHex() bool {
	return '0' <= l.rune && l.rune <= '9' || 'a' <= l.Lower() && l.Lower() <= 'f'
}
