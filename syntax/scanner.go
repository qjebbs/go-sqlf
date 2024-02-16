package syntax

import "strconv"

// scanFn is the lexical scan function
type scanFn func(*scanner) scanFn

// scanner is the lexical scanner
type scanner struct {
	*lexerHelper

	tokens []*token
	token  *token
	state  scanFn
}

func newScanner(input string) *scanner {
	s := &scanner{
		lexerHelper: newLexerHelper(input),
		state:       scanPlain,
	}
	return s
}

func (s *scanner) emitToken(t TokenType, kind litKind, bad bool) {
	s.tokens = append(s.tokens, &token{
		typ:   t,
		kind:  kind,
		bad:   bad,
		start: s.start.offset,
		end:   s.current.offset,
		pos:   s.start.Pos,
		lit:   s.input[s.start.offset:s.current.offset],
	})
}

// NextToken finds the next token
func (s *scanner) NextToken() bool {
	for (len(s.tokens) == 0) && s.state != nil {
		s.state = s.state(s)
	}
	if len(s.tokens) > 0 {
		s.token = s.tokens[0]
		s.tokens = s.tokens[1:]
		return true
	}
	return false
}

func scanPlain(s *scanner) scanFn {
	s.StartToken()
	for r := s.rune; r != EOF; r = s.Next() {
		switch r {
		case '$', '?':
			if s.Peek() == r {
				s.Next()
				continue
			}
			if s.current.offset > s.start.offset {
				s.emitToken(_Plain, _StringLit, false)
			}
			return scanRef
		case '#':
			if s.current.offset > s.start.offset {
				s.emitToken(_Plain, _StringLit, false)
			}
			return scanFunc
		case '\'', '"', '`':
			return scanQuotedPlain
		}
	}
	// EOF
	if s.current.offset > s.start.offset {
		s.emitToken(_Plain, _StringLit, false)
		return scanPlain
	}
	s.emitToken(_EOF, _StringLit, false)
	return nil
}

func scanRef(s *scanner) scanFn {
	s.StartToken()
	s.Next()
	s.emitToken(_Ref, _StringLit, false)
	return scanIndex
}

func scanIndex(s *scanner) scanFn {
	switch s.rune {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.StartToken()
		for r := s.rune; r != EOF; r = s.Next() {
			if r < '0' || r > '9' {
				break
			}
		}
		s.emitToken(_Literal, _NumberLit, false)
	}
	return scanPlain
}

func scanQuotedPlain(s *scanner) scanFn {
	quoter := s.rune
	for r := s.Next(); r != EOF; r = s.Next() {
		if r == quoter {
			if quoter == '\'' && s.Peek() == '\'' {
				s.Next()
				continue
			}
			s.Next()
			s.emitToken(_Plain, _StringLit, false)
			return scanPlain
		}
	}
	// EOF
	s.emitToken(_Plain, _StringLit, true)
	return scanPlain
}

func scanFunc(s *scanner) scanFn {
	s.StartToken()
	s.Next()
	s.emitToken(_Hash, _StringLit, false)
	return scanFuncName
}

func scanFuncName(s *scanner) scanFn {
	s.StartToken()
	for s.IsLetter() {
		s.Next()
	}
	if !s.Advanced() {
		return scanPlain
	}
	s.emitToken(_Name, _StringLit, false)
	s.StartToken()
	for s.IsDecimal() {
		s.Next()
	}
	if s.Advanced() {
		s.emitToken(_Literal, _NumberLit, false)
		return scanPlain
	}
	if s.rune == '(' {
		s.Next()
		s.emitToken(_Lparen, _StringLit, false)
		return scanFuncArgs
	}
	return scanPlain
}

func scanFuncArgs(s *scanner) scanFn {
	s.SkipWhitespace()
	s.StartToken()
	r := s.rune
	for r != EOF {
		switch r {
		case ',':
			s.Next()
			s.emitToken(_Comma, _StringLit, false)
			return scanFuncArgs
		case ')':
			s.Next()
			s.emitToken(_Rparen, _StringLit, false)
			return scanPlain
		case '\'':
			return scanFuncArgQuoted
		default:
			for r != EOF && r != ',' && r != ')' {
				r = s.Next()
			}
			if s.Advanced() {
				fragment := s.input[s.start.offset:s.current.offset]
				if fragment == "true" || fragment == "false" {
					s.emitToken(_Literal, _BoolLit, false)
					return scanFuncArgs
				}
				if fragment == "null" || fragment == "nil" {
					s.emitToken(_Literal, _NilLit, false)
					return scanFuncArgs
				}
				if _, err := strconv.ParseFloat(fragment, 64); err == nil {
					s.emitToken(_Literal, _NumberLit, false)
					return scanFuncArgs
				}
				s.emitToken(_Name, _StringLit, true)
			}
			return scanFuncArgs
		}
	}
	return scanPlain
}

func scanFuncArgQuoted(s *scanner) scanFn {
	quoter := s.rune
	for r := s.Next(); r != EOF; r = s.Next() {
		if r == quoter {
			if quoter == '\'' && s.Peek() == '\'' {
				s.Next()
				continue
			}
			s.Next()
			s.emitToken(_Literal, _StringLit, false)
			return scanFuncArgs
		}
	}
	// EOF
	s.emitToken(_Literal, _StringLit, true)
	return scanPlain
}
