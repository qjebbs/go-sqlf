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
		scanner: newScanner(input, false),
	}
	if err := p.Parse(); err != nil {
		return nil, err
	}
	return p.c, nil
}

type parser struct {
	*scanner

	bindVarIndex int
	bindVarStyle BindStyle
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
	if p.token.bad {
		return nil, p.syntaxError("invalid bind variable: " + p.token.lit)
	}

	t := getBindVarStyle(p.token)
	if err := p.checkVarStyle(t); err != nil {
		return nil, err
	}

	pos := p.token.pos
	lit := p.token.lit[1:]
	switch p.token.kind {
	case _KindRefPositional:
		p.bindVarIndex++
		return &BindVarExpr{
			typ:   t,
			Index: p.bindVarIndex,
			expr:  expr{node{pos}},
		}, nil
	case _KindRefNumbered:
		val, err := strconv.ParseUint(lit, 10, 64)
		if err != nil {
			return nil, p.syntaxError(err.Error())
		}
		return &BindVarExpr{
			typ:   t,
			Index: int(val),
			expr:  expr{node{pos}},
		}, nil
	case _KindRefNamed:
		return &BindVarExpr{
			typ:  t,
			Name: lit,
			expr: expr{node{pos}},
		}, nil
	default:
		return nil, p.syntaxError("unknown bindvar style: " + p.token.lit)
	}
}

func (p *parser) checkVarStyle(t BindStyle) error {
	if p.bindVarStyle == BindStyleUnknown {
		p.bindVarStyle = t
		return nil
	}
	if p.bindVarStyle != t {
		return p.syntaxError("mixed bindvar styles")
	}
	return nil
}

func getBindVarStyle(token *token) BindStyle {
	switch token.kind {
	case _KindRefPositional:
		return BindStyleQuestion
	}
	prefix := token.lit[:1]
	switch token.kind {
	case _KindRefNumbered:
		switch prefix {
		case "$":
			return BindStyleDollarNumbered
		case ":":
			return BindStyleColonNumbered
		case "?":
			return BindStyleQuestionNumbered
		}
	case _KindRefNamed:
		switch prefix {
		case "@":
			return BindStyleAtNamed
		case ":":
			return BindStyleColonNamed
		}
	}
	return BindStyleUnknown
}
