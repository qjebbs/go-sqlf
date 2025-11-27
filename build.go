package sqlf

import (
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

var _ Builder = (*Fragment)(nil)

// BuildQuery builds the fragment as full query.
func (f *Fragment) BuildQuery(style BindStyle) (query string, args []any, err error) {
	ctx := NewContext(style)
	query, err = f.Build(ctx)
	if err != nil {
		return "", nil, err
	}
	args = ctx.Args()
	return query, args, nil
}

// Build builds the fragment with context.
func (f *Fragment) Build(ctx *Context) (string, error) {
	if f == nil {
		return "", nil
	}
	if ctx == nil {
		return "", fmt.Errorf("nil context")
	}
	body, err := build(ctx, f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(body), nil
}

// build builds the fragment
func build(ctx *Context, fragment *Fragment) (string, error) {
	clause, err := syntax.Parse(fragment.Raw)
	if err != nil {
		return "", fmt.Errorf("parse '%s': %w", fragment.Raw, err)
	}
	built, err := buildClause(ctx, fragment, clause)
	if err != nil {
		return "", fmt.Errorf("build '%s': %w", fragment.Raw, err)
	}
	return built, nil
}

// buildClause builds the parsed clause within current context.
func buildClause(ctx *Context, fragment *Fragment, clause *syntax.Clause) (string, error) {
	props := newProperties(fragment.Args...)
	b := new(strings.Builder)
	for _, decl := range clause.ExprList {
		switch expr := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(expr.Text)
		case *syntax.BindVarExpr:
			if expr.Index < 1 || expr.Index > len(props) {
				return "", fmt.Errorf("invalid bind var index %d", expr.Index)
			}
			s, err := props[expr.Index-1].Build(ctx)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		default:
			return "", fmt.Errorf("unknown expression type %T", expr)
		}
	}
	if err := props.checkUsage(); err != nil {
		return "", fmt.Errorf("build '%s': args %w", fragment.Raw, err)
	}
	return b.String(), nil
}
