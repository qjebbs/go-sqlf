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

// checkUsage checks if all properties are used.
func (p properties) checkUsage() error {
	for i, prop := range p {
		if !prop.Used() {
			return fmt.Errorf("#%d unused", i+1)
		}
	}
	return nil
}
