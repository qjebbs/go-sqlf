package syntax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qjebbs/go-sqlf/v4/internal/util"
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

func (p *parser) want(typs ...TokenType) error {
	if !p.got(typs...) {
		return p.syntaxError(
			fmt.Sprintf("syntax error: unexpected %s, want %s", p.token.typ, strings.Join(util.Map(typs, func(t TokenType) string { return string(t) }), " / ")),
		)
	}
	return nil
}

func (p *parser) got(typs ...TokenType) bool {
	p.NextToken()
	for _, t := range typs {
		if p.token.typ == t {
			return true
		}
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
		case _Escape:
			p.c.ExprList = append(p.c.ExprList, &PlainExpr{
				Text: p.token.lit[:1],
				expr: expr{node{p.token.pos}},
			})
		case _Literal:
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
	startToken := p.token
	pos := startToken.pos
	var t bindStyle
	switch startToken.lit {
	case "$":
		t = bindStyleDollarNumbered
		if err := p.checkVarStyle(t); err != nil {
			return nil, err
		}
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
		return &BindVarExpr{
			typ:   t,
			Index: int(val),
			expr:  expr{node{pos}},
		}, nil
	case "?":
		t = bindStyleQuestion
		p.bindVarIndex++
		if err := p.checkVarStyle(t); err != nil {
			return nil, err
		}
		return &BindVarExpr{
			typ:   t,
			Index: p.bindVarIndex,
			expr:  expr{node{pos}},
		}, nil
	default:
		return nil, p.syntaxError("unknown bindvar style: " + startToken.lit)
	}
}

func (p *parser) checkVarStyle(t bindStyle) error {
	if p.bindVarStyle == 0 {
		p.bindVarStyle = t
	}
	if p.bindVarStyle != t {
		return p.syntaxError("mixed bindvar styles")
	}
	return nil
}
