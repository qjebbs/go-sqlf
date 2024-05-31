package sqlf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/v2/syntax"
)

var builtInFuncs = FuncMap{
	"arg":      funcArg,
	"f":        funcFragment,
	"fragment": funcFragment,
	"join":     funcJoin,
}

func funcJoin(ctx *Context, tmpl, separator string, indexes ...int) (string, error) {
	var err error
	var from, to int
	switch len(indexes) {
	case 0:
	case 1:
		from = indexes[0]
	case 2:
		from = indexes[0]
		to = indexes[1]
	default:
		return "", fmt.Errorf(
			"too many args for #join: want 2-4 got %d",
			len(indexes)+2,
		)
	}
	if to > 0 && from > to {
		return "", fmt.Errorf("invalid index range %d to %d", from, to)
	}
	c, err := syntax.Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse join template '%s': %w", tmpl, err)
	}
	b := new(strings.Builder)
	var calls []*syntax.FuncCallExpr
	for i, expr := range c.ExprList {
		fn, ok := expr.(*syntax.FuncExpr)
		if !ok {
			continue
		}
		f, ok := ctx.fn(fn.Name)
		if !ok {
			return "", fmt.Errorf("unknown function #%s", fn.Name)
		}
		if f.JoinCompatibilityError() != nil {
			return "", fmt.Errorf("function #%s is incompatible with #join: %w", fn.Name, f.joinError)
		}
		call := &syntax.FuncCallExpr{
			Name: fn.Name,
		}
		c.ExprList[i] = call
		calls = append(calls, call)
	}
	if len(calls) == 0 {
		return "", fmt.Errorf("no function in join template '%s' (e.g.: #f, not #f1)", tmpl)
	}
	start := from
	if start <= 0 {
		start = 1
	}
	for i := start; ; i++ {
		for _, call := range calls {
			call.Args = []any{i}
		}
		s, err := buildClause(ctx, c)
		if errors.Is(err, ErrInvalidIndex) {
			if (from > 0 && i == from) ||
				(to > 0 && i <= to) {
				// report index error only if explicitly specified
				return "", err
			}
			break
		}
		if err != nil {
			return "", err
		}
		if s != "" {
			if b.Len() > 0 {
				b.WriteString(separator)
			}
			b.WriteString(s)
		}
		if to > 0 && i >= to {
			break
		}
	}
	return b.String(), nil
}

func funcArg(ctx *Context, i int) (string, error) {
	return ctx.fragment().Args.Build(ctx, i)
}

func funcFragment(ctx *Context, i int) (string, error) {
	return ctx.fragment().Builders.Build(ctx, i)
}
