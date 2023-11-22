package sqls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/qjebbs/go-sqls/syntax"
)

// preprocessor is the type of preprocessing functions.
type preprocessor func(ctx *context, args ...string) (string, error)

var builtInFuncs map[string]preprocessor

func init() {
	builtInFuncs = map[string]preprocessor{
		"join":    join,
		"$":       argumentDollar,
		"?":       argumentQuestion,
		"global$": globalArgumentDollar,
		"global?": globalArgumentQuestion,
		"c":       column,
		"col":     column,
		"column":  column,
		"t":       table,
		"table":   table,
		"s":       segment,
		"seg":     segment,
		"segment": segment,
		"b":       builder,
		"builder": builder,
	}
}

func join(ctx *context, args ...string) (string, error) {
	if len(args) != 2 {
		return "", argError("join(tmpl, sep string)", args)
	}
	tmpl, separator := args[0], args[1]
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
		return "", fmt.Errorf("no function in join template '%s' (e.g.: #col, not #col1)", tmpl)
	}
	for i := 0; ; i++ {
		for _, call := range calls {
			call.Args = []string{strconv.Itoa(i + 1)}
		}
		s, err := buildClause(ctx, c)
		if errors.Is(err, ErrInvalidIndex) {
			break
		}
		if err != nil {
			return "", err
		}
		if s != "" {
			if i > 0 {
				b.WriteString(separator)
			}
			b.WriteString(s)
		}
	}
	return b.String(), nil
}

func argumentDollar(ctx *context, args ...string) (string, error) {
	return arg(ctx, syntax.Dollar, false, args...)
}

func argumentQuestion(ctx *context, args ...string) (string, error) {
	return arg(ctx, syntax.Question, false, args...)
}

func globalArgumentDollar(ctx *context, args ...string) (string, error) {
	return arg(ctx, syntax.Dollar, true, args...)
}

func globalArgumentQuestion(ctx *context, args ...string) (string, error) {
	return arg(ctx, syntax.Question, true, args...)
}

func arg(ctx *context, typ syntax.BindVarStyle, global bool, args ...string) (string, error) {
	if ctx.global.BindVarStyle == 0 {
		ctx.global.BindVarStyle = typ
		// ctx.global.FirstBindvar = ctx.Segment.Raw
	}
	// if ctx.global.BindVarStyle != typ {
	// 	return "", fmt.Errorf("mixed bindvar styles between segments '%s' and '%s'", ctx.global.FirstBindvar, ctx.Segment.Raw)
	// }
	if len(args) != 1 {
		switch ctx.global.BindVarStyle {
		case syntax.Dollar:
			return "", argError("$(i int)", args)
		default:
			return "", argError("?(i int)", args)
		}
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid arg index '%s': %w", args[0], err)
	}
	if global {
		return ctx.global.Arg(i)
	}
	return buildArg(ctx, i)
}

func column(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("column(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildColumn(ctx, i)
}

func table(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("tableName(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildTable(ctx, i)
}

func segment(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("segment(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildSegment(ctx, i)
}

func builder(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("builder(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildBuilder(ctx, i)
}

func argError(sig string, args any) error {
	return fmt.Errorf("bad args for #%s: got %v", sig, args)
}
