package sqlf

// property is the interface for properties.
type property interface {
	Builder
	// Used reports if the property is used.
	Used() bool
	// ReportUsed marks current property as used
	ReportUsed()
}

var _ property = (*defaultProperty)(nil)

type defaultProperty struct {
	value Builder
	used  bool
}

// newBuilderProperty returns a new property.
func newBuilderProperty(value Builder) *defaultProperty {
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
	return p.value.BuildFragment(ctx)
}

func newArgProperty(value any) *defaultProperty {
	return newBuilderProperty(&argBuilder{value})
}

var _ Builder = (*argBuilder)(nil)

type argBuilder struct {
	any
}

// BuildFragment implements FragmentBuilder
func (c *argBuilder) BuildFragment(ctx *Context) (query string, err error) {
	built := ctx.CommitArg(c.any)
	return built, nil
}
