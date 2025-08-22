package sqlf

// Property is the interface for properties.
type Property interface {
	Builder
	// Used reports if the property is used.
	Used() bool
	// ReportUsed marks current property as used
	ReportUsed()
}

var _ Property = (*defaultProperty)(nil)

type defaultProperty struct {
	value Builder
	used  bool
}

// newDefaultProperty returns a new property.
func newDefaultProperty(value Builder) *defaultProperty {
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
