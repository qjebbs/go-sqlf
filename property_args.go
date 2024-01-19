package sqlf

import (
	"github.com/qjebbs/go-sqlf/syntax"
)

// ArgsProperty is the args builder
type ArgsProperty struct {
	*property[any]
}

// NewArgsProperty returns a new ArgsBuilder.
func NewArgsProperty(args []any) *ArgsProperty {
	return &ArgsProperty{
		property: newProperty("args", args),
	}
}

// Build builds the arg at index.
func (b *ArgsProperty) Build(ctx *Context, index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if err := b.validateIndex(index); err != nil {
		return "", err
	}
	i := index - 1
	b.used[i] = true
	built := b.cache[i]
	if built == "" || ctx.BindVarStyle == syntax.Question {
		built = ctx.CommitArg(b.items[i], defaultStyle)
		b.cache[i] = built
	}
	return built, nil
}
