package sqlf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/qjebbs/go-sqlf/syntax"
)

// preprocessor is the type of preprocessing functions.
type preprocessor func(ctx *context, args ...string) (string, error)

var builtInFuncs map[string]preprocessor

func init() {
	builtInFuncs = map[string]preprocessor{
		"c":              funcColumn,
		"column":         funcColumn,
		"t":              funcTable,
		"table":          funcTable,
		"fragment":       funcFragment,
		"builder":        funcBuilder,
		"argDollar":      funcArgDollar,
		"argQuestion":    funcArgQuestion,
		"ctxArgDollar":   funcCtxArgDollar,
		"ctxArgQuestion": funcCtxArgQuestion,
		"join":           funcJoin,
	}
}

func funcJoin(ctx *context, args ...string) (string, error) {
	if len(args) < 2 || len(args) > 4 {
		return "", argError("join(tmpl, sep string, from? int, to? int)", args)
	}
	tmpl, separator := args[0], args[1]
	var err error
	var from, to int
	if len(args) > 2 {
		from, err = strconv.Atoi(args[2])
		if err != nil {
			return "", fmt.Errorf("invalid from index '%s': %w", args[2], err)
		}
	}
	if len(args) > 3 {
		to, err = strconv.Atoi(args[3])
		if err != nil {
			return "", fmt.Errorf("invalid to index '%s': %w", args[3], err)
		}
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
			call.Args = []string{strconv.Itoa(i)}
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

func funcArgDollar(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("argDollar(i int)", args)
	}
	return arg(ctx, false, syntax.Dollar, args[0])
}

func funcArgQuestion(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("argQuestion(i int)", args)
	}
	return arg(ctx, false, syntax.Question, args[0])
}

func funcCtxArgDollar(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("ctxArgDollar$(i int)", args)
	}
	return arg(ctx, true, syntax.Dollar, args[0])
}

func funcCtxArgQuestion(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("ctxArgQuestion$(i int)", args)
	}
	return arg(ctx, true, syntax.Question, args[0])
}

func arg(ctx *context, global bool, style syntax.BindVarStyle, arg string) (string, error) {
	i, err := strconv.Atoi(arg)
	if err != nil {
		return "", fmt.Errorf("invalid arg index '%s': %w", arg, err)
	}
	if global {
		return ctx.global.buildArg(i, style)
	}
	return buildArg(ctx, i, style)
}

func funcColumn(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("column(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildColumn(ctx, i)
}

func funcTable(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("table(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildTable(ctx, i)
}

func funcFragment(ctx *context, args ...string) (string, error) {
	if len(args) != 1 {
		return "", argError("fragment(i int)", args)
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid index '%s': %w", args[0], err)
	}
	return buildFragment(ctx, i)
}

func funcBuilder(ctx *context, args ...string) (string, error) {
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
