package sqlf

import (
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// FragmentContext is the FragmentContext for current fragment building.
type FragmentContext struct {
	Global *Context  // global context
	This   *Fragment // current fragment

	argsBuilt      []string // cache of built args
	columnsBuilt   []string // cache of built columns
	fragmentsBuilt []string // cache of built fragments
	buildersBuilt  []string // cache of built builders

	argsUsed      []bool // flags to indicate if an arg is used
	columnsUsed   []bool // flags to indicate if a column is used
	tableUsed     []bool // flag to indicate if a table is used
	fragmentsUsed []bool // flags to indicate if a fragment is used
	builderUsed   []bool // flags to indicate if a builder is used
}

func newFragmentContext(ctx *Context, f *Fragment) *FragmentContext {
	if f == nil {
		return nil
	}
	return &FragmentContext{
		Global:         ctx,
		This:           f,
		argsBuilt:      make([]string, len(f.Args)),
		columnsBuilt:   make([]string, len(f.Columns)),
		tableUsed:      make([]bool, len(f.Tables)),
		fragmentsBuilt: make([]string, len(f.Fragments)),
		buildersBuilt:  make([]string, len(f.Builders)),
		argsUsed:       make([]bool, len(f.Args)),
		columnsUsed:    make([]bool, len(f.Columns)),
		fragmentsUsed:  make([]bool, len(f.Fragments)),
		builderUsed:    make([]bool, len(f.Builders)),
	}
}

// Arg returns the rendered the bindvar at index.
func (c *FragmentContext) Arg(index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if index > len(c.This.Args) {
		return "", fmt.Errorf("%w: bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.This.Args))
	}
	c.Global.onBindArg(defaultStyle)
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.Global.bindVarStyle == syntax.Question {
		*c.Global.ArgStore = append(*c.Global.ArgStore, c.This.Args[i])
		if c.Global.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(*c.Global.ArgStore))
		}
		c.argsBuilt[i] = built
	}
	return built, nil
}

// Column returns the rendered column at index.
func (c *FragmentContext) Column(index int) (string, error) {
	if index > len(c.This.Columns) {
		return "", fmt.Errorf("%w: column index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.This.Columns))
	}
	i := index - 1
	c.columnsUsed[i] = true
	col := c.This.Columns[i]
	built := c.columnsBuilt[i]
	if built == "" || (c.Global.bindVarStyle == syntax.Question && len(col.Args) > 0) {
		b, err := c.buildColumn(col)
		if err != nil {
			return "", err
		}
		c.columnsBuilt[i] = b
		built = b
	}
	return built, nil
}

func (c *FragmentContext) buildColumn(column *TableColumn) (string, error) {
	if column == nil || column.Raw == "" {
		return "", nil
	}
	seg := &Fragment{
		Raw:    column.Raw,
		Args:   column.Args,
		Tables: []Table{column.Table},
	}
	ctxCol := newFragmentContext(c.Global, seg)
	built, err := build(ctxCol)
	if err != nil {
		return "", err
	}
	// don't check usage of tables
	for i := range ctxCol.tableUsed {
		ctxCol.tableUsed[i] = true
	}
	if err := ctxCol.checkUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", column.Raw, err)
	}
	return built, err
}

// Table returns the rendered table at index.
func (c *FragmentContext) Table(index int) (string, error) {
	if index > len(c.This.Tables) {
		return "", fmt.Errorf("%w: table index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.This.Tables))
	}
	c.tableUsed[index-1] = true
	return string(c.This.Tables[index-1]), nil
}

// Fragment returns the rendered fragment at index.
func (c *FragmentContext) Fragment(index int) (string, error) {
	if index > len(c.This.Fragments) {
		return "", fmt.Errorf("%w: fragment index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.This.Fragments))
	}
	i := index - 1
	c.fragmentsUsed[i] = true
	seg := c.This.Fragments[i]
	built := c.fragmentsBuilt[i]
	if built == "" || (c.Global.bindVarStyle == syntax.Question && len(seg.Args) > 0) {
		b, err := seg.BuildContext(c.Global)
		if err != nil {
			return "", err
		}
		c.fragmentsBuilt[i] = b
		built = b
	}
	return built, nil
}

// Builder returns the rendered builder at index.
func (c *FragmentContext) Builder(index int) (string, error) {
	if index > len(c.This.Builders) {
		return "", fmt.Errorf("%w: builder index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.This.Builders))
	}
	i := index - 1
	c.builderUsed[i] = true
	builder := c.This.Builders[i]
	built := c.buildersBuilt[i]
	if built == "" || c.Global.bindVarStyle == syntax.Question {
		b, err := builder.BuildContext(c.Global)
		if err != nil {
			return "", err
		}
		c.buildersBuilt[i] = b
		built = b
	}
	return built, nil
}
func (c *FragmentContext) checkUsage() error {
	if c == nil {
		return nil
	}
	for i, v := range c.argsUsed {
		if !v {
			return fmt.Errorf("arg %d is not used", i+1)
		}
	}
	for i, v := range c.columnsUsed {
		if !v {
			return fmt.Errorf("column %d is not used", i+1)
		}
	}
	for i, v := range c.tableUsed {
		if !v {
			return fmt.Errorf("table %d is not used", i+1)
		}
	}
	for i, v := range c.fragmentsUsed {
		if !v {
			return fmt.Errorf("fragment %d is not used", i+1)
		}
	}
	for i, v := range c.builderUsed {
		if !v {
			return fmt.Errorf("builder %d is not used", i+1)
		}
	}
	return nil
}
