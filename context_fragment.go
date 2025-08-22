package sqlf

import (
	"errors"
	"fmt"
	"strings"
)

// fragment returns the fragment context.
//
// it's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) fragment() (*fragmentContext, bool) {
	return contextValue(c, func(c *Context) (*fragmentContext, bool) {
		return c.frag, c.frag != nil
	})
}

func (c *Context) mustFragment() (*fragmentContext, error) {
	fc, ok := c.fragment()
	if !ok {
		return nil, errors.New("no fragment context")
	}
	return fc, nil
}

// contextWithFragment returns a new context with the fragment.
func contextWithFragment(ctx *Context, f *Fragment) *Context {
	c, _ := contextWith(ctx, func(c *Context) error {
		c.frag = newFragmentContext(f)
		return nil
	})
	return c
}

// fragmentContext is the context for current fragment building.
type fragmentContext struct {
	Fragment *Fragment
	Args     properties
}

func newFragmentContext(f *Fragment) *fragmentContext {
	if f == nil {
		return &fragmentContext{}
	}
	return &fragmentContext{
		Fragment: f,
		Args:     newProperties(f.Args...),
	}
}

// checkUsage checks if all args and properties are used.
func (c *fragmentContext) checkUsage() error {
	if c == nil {
		return nil
	}
	msgs := make([]string, 0, 5)
	if err := c.Args.checkUsage(); err != nil {
		msgs = append(msgs, fmt.Sprintf(
			"args %s", err.Error(),
		))
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}
