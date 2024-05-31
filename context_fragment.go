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

// fragment returns the fragment context.
//
// it's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) fragment() *FragmentContext {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if c.frag != nil {
			return c.frag
		}
	}
	return nil
}

// contextWithFragment returns a new context with the fragment.
func contextWithFragment(ctx *Context, f *Fragment) *Context {
	return &Context{
		parent: ctx,
		frag:   newFragmentContext(f),
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
