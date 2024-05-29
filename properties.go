package sqlf

import "fmt"

// Properties is a list of properties.
type Properties []Property

// Build builds the propery at i.
// it returns ErrInvalidIndex when the i is out of range, which
// is the required behaviour for a custom #func to be compatible
// with #join.
func (p Properties) Build(ctx *Context, i int) (string, error) {
	if i < 1 || i > len(p) {
		return "", fmt.Errorf("%w: %d", ErrInvalidIndex, i)
	}
	return p[i-1].BuildFragment(ctx)
}

// checkUsage checks if all properties are used.
func (p Properties) checkUsage() error {
	for i, prop := range p {
		if !prop.Used() {
			return fmt.Errorf("#%d unused", i+1)
		}
	}
	return nil
}

// NewFragmentProperties creates new properties from FragmentBuilder
func NewFragmentProperties(builders ...FragmentBuilder) Properties {
	r := make(Properties, 0)
	for _, b := range builders {
		r = append(r, newDefaultProperty(b))
	}
	return r
}

// NewArgsProperties  creates new properties from args
func NewArgsProperties(args ...any) Properties {
	r := make(Properties, 0)
	for _, a := range args {
		r = append(r, newDefaultProperty(&arg{a}))
	}
	return r
}
