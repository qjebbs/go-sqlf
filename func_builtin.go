package sqlf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/syntax"
)

var builtInFuncs = FuncMap{
	"c":                 funcColumn,
	"column":            funcColumn,
	"t":                 funcTable,
	"table":             funcTable,
	"fragment":          funcFragment,
	"builder":           funcBuilder,
	"argDollar":         funcArgDollar,
	"argQuestion":       funcArgQuestion,
	"globalArgDollar":   funcGlobalArgDollar,
	"globalArgQuestion": funcGlobalArgQuestion,
	"join":              funcJoin,
}

func funcJoin(ctx *FragmentContext, tmpl, separator string, indexes ...int) (string, error) {
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
		call := &syntax.FuncCallExpr{
			Name: fn.Name,
		}
		c.ExprList[i] = call
		calls = append(calls, call)
	}
	if len(calls) == 0 {
		return "", fmt.Errorf("no function in join template '%s' (e.g.: #c, not #c1)", tmpl)
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
			if from > 0 && to > 0 {
				// index must be valid if from-to explicitly specified
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

func funcArgDollar(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildArg(i, syntax.Dollar)
}

func funcArgQuestion(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildArg(i, syntax.Question)
}

func funcGlobalArgDollar(ctx *FragmentContext, i int) (string, error) {
	return ctx.Global.BuildArg(i, syntax.Dollar)
}

func funcGlobalArgQuestion(ctx *FragmentContext, i int) (string, error) {
	return ctx.Global.BuildArg(i, syntax.Question)
}

func funcColumn(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildColumn(i)
}

func funcTable(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildTable(i)
}

func funcFragment(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildFragment(i)
}

func funcBuilder(ctx *FragmentContext, i int) (string, error) {
	return ctx.BuildBuilder(i)
}
