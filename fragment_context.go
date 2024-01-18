package sqlf

import (
	"strings"

	"github.com/go-oauth2/oauth2/v4/errors"
)

// FragmentContext is the FragmentContext for current fragment building.
type FragmentContext struct {
	Global     *Context  // global context
	Fragment   *Fragment // current fragment
	Properties *properties
}

func newFragmentContext(ctx *Context, f *Fragment) *FragmentContext {
	if f == nil {
		return nil
	}
	return &FragmentContext{
		Global:     ctx,
		Fragment:   f,
		Properties: newProperties(f),
	}
}

// CheckUsage checks if all args, columns, tables, fragments and builders are used.
func (c *FragmentContext) CheckUsage() error {
	if c == nil {
		return nil
	}
	msgs := make([]string, 0, 5)
	if err := c.Properties.Args.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Properties.Columns.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Properties.Tables.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Properties.Fragments.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Properties.Builders.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}
