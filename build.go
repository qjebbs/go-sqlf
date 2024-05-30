package sqlf

import (
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/v2/syntax"
)

var _ FragmentBuilder = (*Fragment)(nil)
var _ QueryBuilder = (*Fragment)(nil)

// BuildQuery builds the fragment as full query.
func (f *Fragment) BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error) {
	ctx := NewContext(bindVarStyle)
	query, err = f.BuildFragment(ctx)
	if err != nil {
		return "", nil, err
	}
	args = ctx.Args()
	return query, args, nil
}

// BuildFragment builds the fragment with context.
func (f *Fragment) BuildFragment(ctx *Context) (string, error) {
	if f == nil {
		return "", nil
	}
	if ctx == nil {
		return "", fmt.Errorf("nil context")
	}
	ctx = ContextWithFragment(ctx, f)
	body, err := build(ctx, f)
	if err != nil {
		return "", err
	}
	if err := ctx.Fragment().checkUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", f.Raw, err)
	}
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
func build(ctx *Context, fragment *Fragment) (string, error) {
	clause, err := syntax.Parse(fragment.Raw)
	if err != nil {
		return "", fmt.Errorf("parse '%s': %w", fragment.Raw, err)
	}
	built, err := buildClause(ctx, clause)
	if err != nil {
		return "", fmt.Errorf("build '%s': %w", fragment.Raw, err)
	}
	return built, nil
}

// buildClause builds the parsed clause within current context.
func buildClause(ctx *Context, clause *syntax.Clause) (string, error) {
	b := new(strings.Builder)
	for _, decl := range clause.ExprList {
		switch expr := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(expr.Text)
		case *syntax.BindVarExpr:
			args := ctx.Fragment().Args
			if expr.Index < 1 || expr.Index > len(args) {
				return "", fmt.Errorf("invalid bind var index %d", expr.Index)
			}
			s, err := args[expr.Index-1].BuildFragment(ctx)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		case *syntax.FuncCallExpr:
			s, err := evalFunction(ctx, expr.Name, expr.Args)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		case *syntax.FuncExpr:
			return "", fmt.Errorf("unexpected function value at %s, forgot to call it?", expr.Pos())
		default:
			return "", fmt.Errorf("unknown expression type %T", expr)
		}
	}
	return b.String(), nil
}
