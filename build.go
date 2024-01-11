package sqlf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Build builds the fragment.
func (f *Fragment) Build() (query string, args []any, err error) {
	args = make([]any, 0)
	query, err = f.BuildContext(NewContext(&args))
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
	if ctx.ArgStore == nil {
		return "", fmt.Errorf("nil arg store (of *[]any)")
	}
	ctxCur := newFragmentContext(ctx, f)
	body, err := build(ctxCur)
	if err != nil {
		return "", err
	}
	if err := ctxCur.checkUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", ctxCur.Fragment.Raw, err)
	}
	// TODO: check usage of global args
	// check inside BuildContext() is not a good idea,
	// because it's not called for every child fragment/builder,
	// when the building is not complete yet.
	// if err := ctx.checkUsage(); err != nil {
	// 	return "", fmt.Errorf("context: %w", err)
	// }
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
func build(ctx *context) (string, error) {
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
func buildClause(ctx *context, clause *syntax.Clause) (string, error) {
	b := new(strings.Builder)
	for _, decl := range clause.ExprList {
		switch expr := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(expr.Text)
		case *syntax.FuncCallExpr:
			fn := builtInFuncs[expr.Name]
			if fn == nil {
				return "", fmt.Errorf("function '%s' is not found", expr.Name)
			}
			s, err := fn(ctx, expr.Args...)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		case *syntax.BindVarExpr:
			s, err := buildArg(ctx, expr.Index, expr.Type)
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

// Arg renders the bindvar at index.
func buildArg(ctx *context, index int, style syntax.BindVarStyle) (string, error) {
	if index > len(ctx.Fragment.Args) {
		return "", fmt.Errorf("%w: bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(ctx.Fragment.Args))
	}
	ctx.global.onBindArg(style)
	i := index - 1
	ctx.ArgsUsed[i] = true
	built := ctx.ArgsBuilt[i]
	if built == "" || ctx.global.bindVarStyle == syntax.Question {
		*ctx.global.ArgStore = append(*ctx.global.ArgStore, ctx.Fragment.Args[i])
		if ctx.global.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(*ctx.global.ArgStore))
		}
		ctx.ArgsBuilt[i] = built
	}
	return built, nil
}

// Column renders the column at index.
func buildColumn(ctx *context, index int) (string, error) {
	if index > len(ctx.Fragment.Columns) {
		return "", fmt.Errorf("%w: column index %d out of range [1,%d]", ErrInvalidIndex, index, len(ctx.Fragment.Columns))
	}
	i := index - 1
	ctx.ColumnsUsed[i] = true
	col := ctx.Fragment.Columns[i]
	built := ctx.ColumnsBuilt[i]
	if built == "" || (ctx.global.bindVarStyle == syntax.Question && len(col.Args) > 0) {
		b, err := buildColumn2(ctx, col)
		if err != nil {
			return "", err
		}
		ctx.ColumnsBuilt[i] = b
		built = b
	}
	return built, nil
}

func buildColumn2(ctx *context, c *TableColumn) (string, error) {
	if c == nil || c.Raw == "" {
		return "", nil
	}
	seg := &Fragment{
		Raw:    c.Raw,
		Args:   c.Args,
		Tables: []Table{c.Table},
	}
	ctx = newFragmentContext(ctx.global, seg)
	built, err := build(ctx)
	if err != nil {
		return "", err
	}
	// don't check usage of tables
	for i := range ctx.TableUsed {
		ctx.TableUsed[i] = true
	}
	if err := ctx.checkUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", ctx.Fragment.Raw, err)
	}
	return built, err
}

func buildTable(ctx *context, index int) (string, error) {
	if index > len(ctx.Fragment.Tables) {
		return "", fmt.Errorf("%w: table index %d out of range [1,%d]", ErrInvalidIndex, index, len(ctx.Fragment.Tables))
	}
	ctx.TableUsed[index-1] = true
	return string(ctx.Fragment.Tables[index-1]), nil
}

func buildFragment(ctx *context, index int) (string, error) {
	if index > len(ctx.Fragment.Fragments) {
		return "", fmt.Errorf("%w: fragment index %d out of range [1,%d]", ErrInvalidIndex, index, len(ctx.Fragment.Fragments))
	}
	i := index - 1
	ctx.FragmentsUsed[i] = true
	seg := ctx.Fragment.Fragments[i]
	built := ctx.FragmentsBuilt[i]
	if built == "" || (ctx.global.bindVarStyle == syntax.Question && len(seg.Args) > 0) {
		b, err := seg.BuildContext(ctx.global)
		if err != nil {
			return "", err
		}
		ctx.FragmentsBuilt[i] = b
		built = b
	}
	return built, nil
}

func buildBuilder(ctx *context, index int) (string, error) {
	if index > len(ctx.Fragment.Builders) {
		return "", fmt.Errorf("%w: builder index %d out of range [1,%d]", ErrInvalidIndex, index, len(ctx.Fragment.Builders))
	}
	i := index - 1
	ctx.BuilderUsed[i] = true
	builder := ctx.Fragment.Builders[i]
	built := ctx.BuildersBuilt[i]
	if built == "" || ctx.global.bindVarStyle == syntax.Question {
		b, err := builder.BuildContext(ctx.global)
		if err != nil {
			return "", err
		}
		ctx.BuildersBuilt[i] = b
		built = b
	}
	return built, nil
}
