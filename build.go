package sqlf

import (
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Build builds the fragment.
func (f *Fragment) Build() (query string, args []any, err error) {
	ctx := NewContext()
	query, err = f.BuildContext(ctx)
	if err != nil {
		return "", nil, err
	}
	args, err = ctx.BuiltArgs()
	if err != nil {
		return "", nil, err
	}
	return query, args, nil
}

// BuildContext builds the fragment with context.
func (f *Fragment) BuildContext(ctx *Context) (string, error) {
	if f == nil {
		return "", nil
	}
	if ctx == nil {
		return "", fmt.Errorf("nil context")
	}
	if ctx.argStore == nil {
		return "", fmt.Errorf("nil arg store (of *[]any)")
	}
	ctxCur := newFragmentContext(ctx, f)
	body, err := build(ctxCur)
	if err != nil {
		return "", err
	}
	if err := ctxCur.CheckUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", ctxCur.Fragment.Raw, err)
	}
	// don't do usage check for global context here,
	// because it's called from children fragments/builders,
	// when the building is not complete yet.

	body = strings.TrimSpace(body)
	if body == "" {
		return "", nil
	}
	header, footer := "", ""
	if f.Prefix != "" {
		header = f.Prefix + " "
	}
	if f.Suffix != "" {
		footer = " " + f.Suffix
	}
	return header + body + footer, nil
}

// build builds the fragment
func build(ctx *FragmentContext) (string, error) {
	clause, err := syntax.Parse(ctx.Fragment.Raw)
	if err != nil {
		return "", fmt.Errorf("parse '%s': %w", ctx.Fragment.Raw, err)
	}
	built, err := buildClause(ctx, clause)
	if err != nil {
		return "", fmt.Errorf("build '%s': %w", ctx.Fragment.Raw, err)
	}
	return built, nil
}

// buildClause builds the parsed clause within current context, not updating the ctx.current.
func buildClause(ctx *FragmentContext, clause *syntax.Clause) (string, error) {
	b := new(strings.Builder)
	for _, decl := range clause.ExprList {
		switch expr := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(expr.Text)
		case *syntax.FuncCallExpr:
			s, err := evalFunction(ctx, expr.Name, expr.Args)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		case *syntax.BindVarExpr:
			s, err := ctx.BuildArg(expr.Index, expr.Type)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		default:
			return "", fmt.Errorf("unknown expression type %T", expr)
		}
	}
	return b.String(), nil
}
