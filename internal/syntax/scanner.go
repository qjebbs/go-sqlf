package syntax

// scanFn is the lexical scan function
type scanFn func(*scanner) scanFn

// scanner is the lexical scanner
type scanner struct {
	*lexerHelper

	tokens []*token
	token  *token
	state  scanFn

	interpolating bool
}

func newScanner(input string, interpolating bool) *scanner {
	s := &scanner{
		lexerHelper:   newLexerHelper(input),
		state:         scanPlain,
		interpolating: interpolating,
	}
	return s
}

func (s *scanner) emitToken(t TokenType, kind kind, bad bool) {
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

// Rewind rewinds the scanner to the previous token
func (s *scanner) Rewind(current, rewinds *token) {
	s.tokens = append([]*token{rewinds}, s.tokens...)
	s.token = current
}

func scanPlain(s *scanner) scanFn {
	s.StartToken()
	defer func() {
		if s.current.offset > s.start.offset {
			s.emitToken(_Raw, 0, false)
		}
	}()

	var bindExtra1, bindExtra2 = '?', '?'
	if s.interpolating {
		bindExtra1 = ':'
		bindExtra2 = '@'
	}

	// In interpolating mode, we intentionally do not support non-standard quote
	// styles [name]. This is because a simple lexer cannot reliably distinguish
	// a quoted identifier (e.g., [column]) from other SQL constructs like array
	// access (e.g., array_column[1]). Attempting to do so without a full parser
	// would be complex and error-prone.
	//
	// The trade-off is that text within an identifier that resembles a bind
	// variable (e.g., [my-col-$1]) may be incorrectly interpreted as a bind
	// variable, breaking the identifier and interpolating. This is an acceptable
	// limitation to avoid the complexity of a full SQL parser.

	for r := s.rune; r != EOF; r = s.Next() {
		switch r {
		case '$', '?', bindExtra1, bindExtra2:
			if s.Peek() == r {
				if s.interpolating {
					// In interpolating mode, $$, ??, ::, @@ are not escapes,
					// neither are they bind vars. Just treat them as plain text.
					s.Next()
					continue
				}
				return scanEscape
			}
			return scanRef
		case '\'', '"', '`':
			return scanQuoted
		}
	}
	return scanEnd
}

func scanEscape(s *scanner) scanFn {
	s.StartToken()
	s.Next()
	s.Next()
	s.emitToken(_Escape, 0, false)
	return scanPlain
}

func scanRef(s *scanner) scanFn {
	s.StartToken()
	prefix := s.rune
	s.Next()

	switch prefix {
	case '?':
		if !s.interpolating || !s.IsDigit() {
			s.emitToken(_BindVar, _KindBindVarPositional, false)
			return scanPlain
		}
		for s.IsDigit() {
			s.Next()
		}
		s.emitToken(_BindVar, _KindBindVarNumbered, false)
		return scanPlain
	case '$':
		if !s.IsDigit() {
			if s.interpolating {
				s.emitToken(_Raw, 0, true)
			} else {
				s.emitToken(_BindVar, _KindBindVarNumbered, true)
			}
			return scanPlain
		}
		for s.IsDigit() {
			s.Next()
		}
		s.emitToken(_BindVar, _KindBindVarNumbered, false)
		return scanPlain
	case '@':
		if !s.IsLetter() {
			s.emitToken(_BindVar, _KindBindVarNamed, true)
			return scanPlain
		}
		for s.IsLetter() || s.IsDigit() {
			s.Next()
		}
		s.emitToken(_BindVar, _KindBindVarNamed, false)
		return scanPlain
	case ':':
		if s.IsLetter() {
			for s.IsLetter() || s.IsDigit() {
				s.Next()
			}
			s.emitToken(_BindVar, _KindBindVarNamed, false)
			return scanPlain
		}
		if s.IsDigit() {
			for s.rune >= '0' && s.rune <= '9' {
				s.Next()
			}
			s.emitToken(_BindVar, _KindBindVarNumbered, false)
			return scanPlain
		}
		s.emitToken(_BindVar, _KindBindVarNumbered, true)
		return scanPlain
	default:
		return scanPlain
	}
}

func scanQuoted(s *scanner) scanFn {
	s.StartToken()
	quoter := s.rune
	for r := s.Next(); r != EOF; r = s.Next() {
		if r == quoter {
			if s.Peek() == quoter {
				s.Next()
				continue
			}
			s.Next()
			s.emitToken(_Literal, _KindLitString, false)
			return scanPlain
		}
	}
	// EOF
	s.emitToken(_Literal, _KindLitString, true)
	return scanEnd
}

func scanEnd(s *scanner) scanFn {
	s.StartToken()
	s.emitToken(_EOF, 0, false)
	return nil
}
