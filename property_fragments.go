package sqlf

import (
	"github.com/qjebbs/go-sqlf/syntax"
)

// FragmentsProperty is the fragments property
type FragmentsProperty struct {
	*propertyBase[*Fragment]
}

// NewFragmentsProperty returns a new FragmentsProperty.
func NewFragmentsProperty(fragments []*Fragment) *FragmentsProperty {
	return &FragmentsProperty{
		propertyBase: newPropertyBase("fragments", fragments),
	}
}

// Build builds the arg at index, with cache.
func (b *FragmentsProperty) Build(ctx *Context, index int) (string, error) {
	if err := b.validateIndex(index); err != nil {
		return "", err
	}
	i := index - 1
	b.used[i] = true
	fragment := b.items[i]
	built := b.cache[i]
	if built == "" || (ctx.bindVarStyle == syntax.Question && len(fragment.Args) > 0) {
		r, err := fragment.BuildContext(ctx)
		if err != nil {
			return "", err
		}
		b.cache[i] = r
		built = r
	}
	return built, nil
}
