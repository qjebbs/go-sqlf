package syntax

import (
	"strconv"
)

// ParseForInterpolating parses the input for interpolation and returns the list of expressions.
// It always returns a Clause, even if there are errors during parsing.
func ParseForInterpolating(input string) (*Clause, error) {
	p := &inerpolatingParser{
		parser: &parser{
			bindVarStyle: BindStyleUnknown,
			scanner:      newScanner(input, true),
		},
	}
	p.Parse()
	return p.c, p.err
}

type inerpolatingParser struct {
	*parser
	err error
}

func (p *inerpolatingParser) Parse() {
	p.c = &Clause{}
L:
	for p.NextToken() {
		switch p.token.typ {
		case _EOF:
			break L
		case _Ref:
			exprs, err := p.bindVarExprInterpolating()
			if err != nil && p.err == nil {
				p.err = err
			}
			p.c.ExprList = append(p.c.ExprList, exprs...)
		default:
			p.c.ExprList = append(p.c.ExprList, &PlainExpr{
				Text: p.token.lit,
				expr: expr{node{p.token.pos}},
			})
		}
	}
}

func (p *inerpolatingParser) bindVarExprInterpolating() (exprs []Expr, err error) {
	startToken := p.token
	pos := startToken.pos

	retrievedTokens := []*token{startToken}
	defer func() {
		if exprs == nil {
			exprs = make([]Expr, 0, len(retrievedTokens))
			for _, tk := range retrievedTokens {
				exprs = append(exprs, &PlainExpr{
					Text: tk.lit,
					expr: expr{node{tk.pos}},
				})
			}
		}
	}()

	var t BindStyle
	switch startToken.lit {
	case "$":
		t = BindStyleDollarNumbered
		if err := p.checkVarStyle(t); err != nil {
			return nil, err
		}
		err = p.want(_Literal)
		retrievedTokens = append(retrievedTokens, p.token)
		if err != nil {
			// not report error when it's not a valid bindvar, parsed as plain text.
			// same for all other cases below
			return nil, nil
		}
		if p.token.kind != _NumberLit {
			return nil, nil
		}
		val, err := strconv.ParseUint(p.token.lit, 10, 64)
		if err != nil {
			return nil, err
		}
		return []Expr{&BindVarExpr{
			typ:   t,
			Index: int(val),
			expr:  expr{node{pos}},
		}}, nil
	case "?":
		err = p.want(_Literal, _Plain, _EOF)
		retrievedTokens = append(retrievedTokens, p.token)
		if err != nil {
			return nil, nil
		}
		var index int
		if p.token.typ == _Literal {
			t = BindStyleQuestionNumbered
			if err := p.checkVarStyle(t); err != nil {
				return nil, err
			}
			if p.token.kind != _NumberLit {
				return nil, nil
			}
			val, err := strconv.ParseUint(p.token.lit, 10, 64)
			if err != nil {
				return nil, err
			}
			index = int(val)
		} else {
			// plain/EOF following "?" means it's not a numbered bindvar
			t = BindStyleQuestion
			if err := p.checkVarStyle(t); err != nil {
				return nil, err
			}
			p.bindVarIndex++
			index = p.bindVarIndex
			p.Rewind(startToken, p.token)
			retrievedTokens = retrievedTokens[:len(retrievedTokens)-1]
		}
		return []Expr{&BindVarExpr{
			typ:   t,
			Index: index,
			expr:  expr{node{pos}},
		}}, nil
	case ":":
		err = p.want(_Literal, _Name)
		retrievedTokens = append(retrievedTokens, p.token)
		if err != nil {
			return nil, nil
		}
		if p.token.typ == _Name {
			t = BindStyleColonNamed
			if err := p.checkVarStyle(t); err != nil {
				return nil, err
			}
			return []Expr{&BindVarExpr{
				typ:  t,
				Name: p.token.lit,
				expr: expr{node{pos}},
			}}, nil
		}
		t = BindStyleColonNumbered
		if err := p.checkVarStyle(t); err != nil {
			return nil, err
		}
		if p.token.kind != _NumberLit {
			return nil, nil
		}
		val, err := strconv.ParseUint(p.token.lit, 10, 64)
		if err != nil {
			return nil, err
		}
		return []Expr{&BindVarExpr{
			typ:   t,
			Index: int(val),
			expr:  expr{node{pos}},
		}}, nil
	case "@":
		err = p.want(_Name)
		retrievedTokens = append(retrievedTokens, p.token)
		if err != nil {
			return nil, nil
		}
		t = BindStyleAtNamed
		if err := p.checkVarStyle(t); err != nil {
			return nil, err
		}
		return []Expr{&BindVarExpr{
			typ:  t,
			Name: p.token.lit,
			expr: expr{node{pos}},
		}}, nil
	}
	return nil, nil
}
