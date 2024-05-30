package sqlf

import (
	"errors"
	"fmt"
	"strings"
)

// FragmentContext is the context for current fragment building.
type FragmentContext struct {
	Fragment *Fragment
	Args     Properties
	Builders Properties
}

// Fragment returns the fragment context.
func (c *Context) Fragment() *FragmentContext {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if c.fragment != nil {
			return c.fragment
		}
	}
	return nil
}

// ContextWithFragment returns a new context with the fragment.
func ContextWithFragment(ctx *Context, f *Fragment) *Context {
	return &Context{
		parent:   ctx,
		fragment: newFragmentContext(f),
	}

}

func newFragmentContext(f *Fragment) *FragmentContext {
	if f == nil {
		return &FragmentContext{}
	}
	return &FragmentContext{
		Fragment: f,
		Args:     NewArgsProperties(f.Args...),
		Builders: NewFragmentProperties(f.Fragments...),
	}
}

// checkUsage checks if all args and properties are used.
func (c *FragmentContext) checkUsage() error {
	if c == nil {
		return nil
	}
	msgs := make([]string, 0, 5)
	if err := c.Args.checkUsage(); err != nil {
		msgs = append(msgs, fmt.Sprintf(
			"args %s", err.Error(),
		))
	}
	if err := c.Builders.checkUsage(); err != nil {
		msgs = append(msgs, fmt.Sprintf(
			"properties %s", err.Error(),
		))
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}
