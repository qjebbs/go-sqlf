package sqlf

import "github.com/qjebbs/go-sqlf/v2/syntax"

// Property is the interface for properties.
type Property interface {
	FragmentBuilder
	// Used reports if the property is used.
	Used() bool
	// ReportUsed marks current property as used
	ReportUsed()
}

var _ Property = (*defaultProperty)(nil)

type defaultProperty struct {
	value FragmentBuilder
	used  bool
	cache string
}

// newDefaultProperty returns a new property.
func newDefaultProperty(value FragmentBuilder) *defaultProperty {
	return &defaultProperty{
		value: value,
	}
}

// ReportUsed reports the item is used
func (p *defaultProperty) ReportUsed() {
	p.used = true
}

// Used returns true if the column is used.
func (p *defaultProperty) Used() bool {
	return p.used
}

// BuildFragment builds the fragment.
func (p *defaultProperty) BuildFragment(ctx *Context) (string, error) {
	p.used = true
	if p.cache == "" || ctx.BindVarStyle() == syntax.Question {
		r, err := p.value.BuildFragment(ctx)
		if err != nil {
			return "", err
		}
		p.cache = r
	}
	return p.cache, nil
}
