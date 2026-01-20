package sqlf

import (
	"fmt"
	"strings"

	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

var _ Builder = (*Fragment)(nil)

// Build builds f into a query string and its corresponding arguments.
// It creates a new context to ensure that the original context is not modified.
func (f *Fragment) Build(ctx *Context) (query string, args []any, err error) {
	return Build(ctx, f)
}

// BuildTo builds the SQL fragment for f based on the given
// context, and commits any arguments to the context.
func (f *Fragment) BuildTo(ctx *Context) (string, error) {
	if f == nil {
		return "", nil
	}
	if ctx == nil {
		return "", fmt.Errorf("nil context")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	body, err := build(ctx, f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(body), nil
}

// build builds the fragment
func build(ctx *Context, fragment *Fragment) (string, error) {
	clause, err := syntax.Parse(fragment.raw)
	if err != nil {
		return "", fmt.Errorf("parse '%s': %w", fragment.raw, err)
	}
	built, err := buildClause(ctx, fragment, clause)
	if err != nil {
		return "", fmt.Errorf("build '%s': %w", fragment.raw, err)
	}
	return built, nil
}

// buildClause builds the parsed clause within current context.
func buildClause(ctx *Context, fragment *Fragment, clause *syntax.Clause) (string, error) {
	props := newProperties(fragment.args...)
	b := new(strings.Builder)
	for _, decl := range clause.ExprList {
		switch expr := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(expr.Text)
		case *syntax.IdentityExpr:
			b.WriteString(ctx.Dialect().QuoteStyle().QuoteIdentifier(expr.Name))
		case *syntax.BindVarExpr:
			if expr.Name != "" {
				// should not happen
				return "", fmt.Errorf("named bind vars not supported by *Fragment")
			}
			if expr.Index < 1 || expr.Index > len(props) {
				return "", fmt.Errorf("invalid bind var index %d", expr.Index)
			}
			s, err := props[expr.Index-1].BuildTo(ctx)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		default:
			return "", fmt.Errorf("unknown expression type %T", expr)
		}
	}
	if !fragment.noUsageCheck {
		if err := props.checkUsage(); err != nil {
			return "", fmt.Errorf("build '%s': args %w", fragment.raw, err)
		}
	}
	return b.String(), nil
}
