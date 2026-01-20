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
	bindVarStyle BindVarStyle
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
			expr, err := p.bindVarExpr()
			if err != nil {
				return err
			}
			p.c.ExprList = append(p.c.ExprList, expr)
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
			expr, err := p.literalExpr()
			if err != nil {
				return err
			}
			p.c.ExprList = append(p.c.ExprList, expr)
		default:
			return p.syntaxError("unexpected token " + string(p.token.typ))
		}
	}
	return nil
}

func (p *parser) literalExpr() (Expr, error) {
	if p.token.bad {
		return nil, p.syntaxError("invalid literal: " + p.token.lit)
	}

	pos := p.token.pos
	lit := p.token.lit
	switch p.token.kind {
	case _KindLitNumber:
		return &PlainExpr{
			Text: lit,
			expr: expr{node{pos}},
		}, nil
	case _KindLitString:
		quoter := lit[:1]
		if quoter == "'" {
			return &PlainExpr{
				Text: lit,
				expr: expr{node{pos}},
			}, nil
		}
		unquoted, err := unquoteString(lit)
		if err != nil {
			return nil, p.syntaxError(err.Error())
		}
		return &IdentityExpr{
			Name: unquoted,
			expr: expr{node{pos}},
		}, nil
	default:
		return nil, p.syntaxError("unknown literal type: " + p.token.lit)
	}
}

func unquoteString(lit string) (string, error) {
	if len(lit) < 2 {
		return "", fmt.Errorf("invalid string literal: %s", lit)
	}
	quote := lit[0]
	if lit[len(lit)-1] != quote {
		return "", fmt.Errorf("mismatched quotes in string literal: %s", lit)
	}
	content := lit[1 : len(lit)-1]
	switch quote {
	case '\'':
		return strings.ReplaceAll(content, "''", "'"), nil
	case '"':
		return strings.ReplaceAll(content, `""`, `"`), nil
	case '`':
		return strings.ReplaceAll(content, "``", "`"), nil
	// case '[':
	// 	return strings.ReplaceAll(content, "]]", "]"), nil
	default:
		return "", fmt.Errorf("unknown quote character: %c", quote)
	}
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

func (p *parser) checkVarStyle(t BindVarStyle) error {
	if p.bindVarStyle == BindVarStyleUnknown {
		p.bindVarStyle = t
		return nil
	}
	if p.bindVarStyle != t {
		return p.syntaxError("mixed bindvar styles")
	}
	return nil
}

func getBindVarStyle(token *token) BindVarStyle {
	switch token.kind {
	case _KindRefPositional:
		return BindVarStyleQuestion
	}
	prefix := token.lit[:1]
	switch token.kind {
	case _KindRefNumbered:
		switch prefix {
		case "$":
			return BindVarStyleDollarNumbered
		case ":":
			return BindVarStyleColonNumbered
		case "?":
			return BindVarStyleQuestionNumbered
		}
	case _KindRefNamed:
		switch prefix {
		case "@":
			return BindVarStyleAtNamed
		case ":":
			return BindVarStyleColonNamed
		}
	}
	return BindVarStyleUnknown
}
