package sqlf

import "fmt"

// properties is a list of properties.
type properties []property

// newProperties  creates new properties from args.
// It's useful for creating global arg properties shared between fragments,
// see examples for ContextWithFuncs() for how to use it.
func newProperties(args ...any) properties {
	r := make(properties, 0)
	for _, a := range args {
		if f, ok := a.(Builder); ok {
			r = append(r, newBuilderProperty(f))
		} else {
			r = append(r, newArgProperty(a))
		}
	}
	return r
}

// Build builds the propery at i.
// it returns ErrInvalidIndex when the i is out of range, which
// is the required behaviour for a custom #func to be compatible
// with #join.
//
// See examples for ContextWithFuncs() for how to use it.
func (p properties) Build(ctx *Context, i int) (string, error) {
	if i < 1 || i > len(p) {
		return "", fmt.Errorf("%w: %d", ErrInvalidIndex, i)
	}
	return p[i-1].BuildFragment(ctx)
}

// checkUsage checks if all properties are used.
func (p properties) checkUsage() error {
	for i, prop := range p {
		if !prop.Used() {
			return fmt.Errorf("#%d unused", i+1)
		}
	}
	return nil
}
