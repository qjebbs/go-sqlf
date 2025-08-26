package syntax

import (
	"fmt"
	"strconv"
)

// Parse parses the input and returns the list of expressions.
func Parse(input string) (*Clause, error) {
	p := &parser{
		scanner: newScanner(input),
	}
	if err := p.Parse(); err != nil {
		return nil, err
	}
	return p.c, nil
}

type parser struct {
	*scanner

	bindVarIndex int
	bindVarStyle bindStyle
	// buf []token

	c *Clause
}

func (p *parser) want(t TokenType) error {
	if !p.got(t) {
		return p.syntaxError(
			fmt.Sprintf("syntax error: unexpected %s, want %s", p.token.typ, t),
		)
	}
	return nil
}

func (p *parser) got(tok TokenType) bool {
	p.NextToken()
	if p.token.typ == tok {
		return true
	}
	return false
}

func (p *parser) syntaxError(msg string) error {
	return fmt.Errorf("%d: syntax error: %s", p.token.start, msg)
}

func (p *parser) Parse() error {
	p.c = &Clause{}
L:
	for p.NextToken() {
		switch p.token.typ {
		case _EOF:
			break L
		case _Ref:
			d, err := p.bindVarExpr()
			if err != nil {
				return err
			}
			p.c.ExprList = append(p.c.ExprList, d)
		case _Plain:
			p.c.ExprList = append(p.c.ExprList, &PlainExpr{
				Text: p.token.lit,
				expr: expr{node{p.token.pos}},
			})
		default:
			return p.syntaxError("unexpected token " + string(p.token.typ))
		}
	}
	return nil
}

func (p *parser) bindVarExpr() (Expr, error) {
	pos := p.token.pos
	var t bindStyle
	switch p.token.lit {
	case "$":
		t = bindStyleDollar
		p.bindVarIndex++
		if p.bindVarStyle == 0 {
			p.bindVarStyle = t
		}
		if p.bindVarStyle != t {
			return nil, p.syntaxError("mixed bindvar styles")
		}
	case "?":
		t = bindStyleQuestion
		p.bindVarIndex++
		if p.bindVarStyle == 0 {
			p.bindVarStyle = t
		}
		if p.bindVarStyle != t {
			return nil, p.syntaxError("mixed bindvar styles")
		}
	}
	index := p.bindVarIndex
	if t != bindStyleQuestion {
		if err := p.want(_Literal); err != nil {
			return nil, err
		}
		if p.token.kind != _NumberLit {
			return nil, p.syntaxError("unexpected '" + p.token.lit + "', want bindvar index")
		}
		val, err := strconv.ParseUint(p.token.lit, 10, 64)
		if err != nil {
			return nil, p.syntaxError(err.Error())
		}
		index = int(val)
	}
	return &BindVarExpr{
		typ:   t,
		Index: index,
		expr:  expr{node{pos}},
	}, nil
}
