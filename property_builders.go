package sqlf

import (
	"github.com/qjebbs/go-sqlf/syntax"
)

// BuildersProperty is the Builders property
type BuildersProperty struct {
	*property[Builder]
}

// NewBuildersProperty returns a new BuildersProperty.
func NewBuildersProperty(builders ...Builder) *BuildersProperty {
	return &BuildersProperty{
		property: newProperty("builders", builders),
	}
}

// Build builds the builder at index.
func (b *BuildersProperty) Build(ctx *Context, index int) (string, error) {
	if err := b.validateIndex(index); err != nil {
		return "", err
	}
	i := index - 1
	b.used[i] = true
	builder := b.items[i]
	built := b.cache[i]
	if built == "" || ctx.BindVarStyle == syntax.Question {
		r, err := builder.BuildContext(ctx)
		if err != nil {
			return "", err
		}
		b.cache[i] = r
		built = r
	}
	return built, nil
}
