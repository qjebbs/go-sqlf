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

// Rewind rewinds the scanner to the previous token
func (s *scanner) Rewind(current, rewinds *token) {
	s.tokens = append([]*token{rewinds}, s.tokens...)
	s.token = current
}

func scanPlain(s *scanner) scanFn {
	s.StartToken()
	var extra1, extra2 = '?', '?'
	if s.interpolating {
		extra1 = ':'
		extra2 = '@'
	}
	for r := s.rune; r != EOF; r = s.Next() {
		switch r {
		case '$', '?', extra1, extra2:
			if !s.interpolating && s.Peek() == r {
				if s.current.offset > s.start.offset {
					s.emitToken(_Plain, _StringLit, false)
				}
				return scanEscape
			}
			if s.current.offset > s.start.offset {
				s.emitToken(_Plain, _StringLit, false)
			}
			return scanRef
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

func scanEscape(s *scanner) scanFn {
	s.StartToken()
	s.Next()
	s.Next()
	s.emitToken(_Escape, _StringLit, false)
	return scanPlain
}

func scanRef(s *scanner) scanFn {
	s.StartToken()
	s.Next()
	s.emitToken(_Ref, _StringLit, false)
	switch s.rune {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return scanIndex
	}
	return scanName
}

func scanIndex(s *scanner) scanFn {
	s.StartToken()
	for r := s.rune; r != EOF; r = s.Next() {
		if r < '0' || r > '9' {
			break
		}
	}
	s.emitToken(_Literal, _NumberLit, false)
	return scanPlain
}

func scanName(s *scanner) scanFn {
	s.StartToken()
	for s.IsLetter() {
		s.Next()
	}
	if s.Advanced() {
		s.emitToken(_Name, _StringLit, false)
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
