package sqlf

import (
	"errors"
	"fmt"
	"strings"
)

// FragmentContext is the FragmentContext for current fragment building.
type FragmentContext struct {
	Global *Context // global context

	Raw      string
	Args     Properties
	Builders Properties
}

func newFragmentContext(ctx *Context, f *Fragment) *FragmentContext {
	if f == nil {
		return nil
	}
	return &FragmentContext{
		Global:   ctx,
		Raw:      f.Raw,
		Args:     NewArgsProperties(f.Args...),
		Builders: NewFragmentProperties(f.Fragments...),
	}
}

// CheckUsage checks if all args and properties are used.
func (c *FragmentContext) CheckUsage() error {
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
