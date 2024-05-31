package sqlf

import "fmt"

// Properties is a list of properties.
type Properties []Property

// Build builds the propery at i.
// it returns ErrInvalidIndex when the i is out of range, which
// is the required behaviour for a custom #func to be compatible
// with #join.
//
// See examples for ContextWithFuncs() for how to use it.
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

// NewFragmentProperties creates new properties from FragmentBuilder.
// It's useful for creating global fragment properties shared between fragments,
// see examples for ContextWithFuncs() for how to use it.
func NewFragmentProperties(fragments ...FragmentBuilder) Properties {
	r := make(Properties, 0)
	for _, f := range fragments {
		r = append(r, newDefaultProperty(f))
	}
	return r
}

// NewArgsProperties  creates new properties from args.
// It's useful for creating global arg properties shared between fragments,
// see examples for ContextWithFuncs() for how to use it.
func NewArgsProperties(args ...any) Properties {
	r := make(Properties, 0)
	for _, a := range args {
		r = append(r, newDefaultProperty(&arg{a}))
	}
	return r
}

var _ FragmentBuilder = (*arg)(nil)

type arg struct {
	any
}

// BuildFragment implements FragmentBuilder
func (c *arg) BuildFragment(ctx *Context) (query string, err error) {
	built := ctx.CommitArg(c.any)
	return built, nil
}
