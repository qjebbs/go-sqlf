package syntax

// ParseForInterpolating parses the input for interpolation and returns the list of expressions.
// It always returns a Clause, even if there are errors during parsing.
func ParseForInterpolating(input string, bindStyle BindVarStyle) (*Clause, error) {
	p := &inerpolatingParser{
		parser: &parser{
			bindVarStyle: bindStyle,
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
			expr, err := p.bindVarExpr()
			if err != nil && p.err == nil {
				p.err = err
			}
			p.c.ExprList = append(p.c.ExprList, expr)
		default:
			p.c.ExprList = append(p.c.ExprList, &PlainExpr{
				Text: p.token.lit,
				expr: expr{node{p.token.pos}},
			})
		}
	}
}

func (p *inerpolatingParser) bindVarExpr() (exp Expr, err error) {
	defer func() {
		if exp == nil {
			exp = &PlainExpr{
				Text: p.token.lit,
				expr: expr{node{p.token.pos}},
			}
		}
	}()
	return p.parser.bindVarExpr()
}
