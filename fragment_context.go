package sqlf

import (
	"errors"
	"strings"
)

// FragmentContext is the FragmentContext for current fragment building.
type FragmentContext struct {
	Global *Context // global context

	Raw       string
	Args      *ArgsProperty
	Columns   *ColumnsProperty
	Tables    *TablesProperty
	Fragments *FragmentsProperty
	Builders  *BuildersProperty
}

func newFragmentContext(ctx *Context, f *Fragment) *FragmentContext {
	if f == nil {
		return nil
	}
	return &FragmentContext{
		Global:    ctx,
		Raw:       f.Raw,
		Args:      NewArgsProperty(f.Args...),
		Columns:   NewColumnsProperty(f.Columns...),
		Tables:    NewTablesProperty(f.Tables...),
		Fragments: NewFragmentsProperty(f.Fragments...),
		Builders:  NewBuildersProperty(f.Builders...),
	}
}

// CheckUsage checks if all properties are used.
func (c *FragmentContext) CheckUsage() error {
	if c == nil {
		return nil
	}
	msgs := make([]string, 0, 5)
	if err := c.Args.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Columns.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Tables.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Fragments.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if err := c.Builders.CheckUsage(); err != nil {
		msgs = append(msgs, err.Error())
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}
