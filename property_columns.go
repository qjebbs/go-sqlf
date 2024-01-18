package sqlf

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/syntax"
)

// ColumnsProperty is the columns property
type ColumnsProperty struct {
	*propertyBase[*Column]
}

// NewColumnsProperty returns a new ColumnsProperty.
func NewColumnsProperty(columns []*Column) *ColumnsProperty {
	return &ColumnsProperty{
		propertyBase: newPropertyBase("columns", columns),
	}
}

// Build builds the arg at index, with cache.
func (b *ColumnsProperty) Build(ctx *Context, index int) (string, error) {
	if err := b.validateIndex(index); err != nil {
		return "", err
	}
	i := index - 1
	b.used[i] = true
	built := b.cache[i]
	if built == "" || (ctx.bindVarStyle == syntax.Question && len(b.items[i].Args) > 0) {
		r, err := b.buildColumn(ctx, b.items[i])
		if err != nil {
			return "", err
		}
		b.cache[i] = r
		built = r
	}
	return built, nil
}

func (b *ColumnsProperty) buildColumn(ctx *Context, column *Column) (string, error) {
	if column == nil || column.Raw == "" {
		return "", nil
	}
	fragment := &Fragment{
		Raw:    column.Raw,
		Args:   column.Args,
		Tables: []Table{column.Table},
	}
	ctxColumn := newFragmentContext(ctx, fragment)
	built, err := build(ctxColumn)
	if err != nil {
		return "", err
	}
	// don't check usage of tables
	ctxColumn.Properties.Tables.ReportUsed(1)
	if err := ctxColumn.CheckUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", column.Raw, err)
	}
	return built, err
}
