package sqlf

import (
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// FragmentContext is the FragmentContext for current fragment building.
type FragmentContext struct {
	Global   *Context  // global context
	Fragment *Fragment // current fragment

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
		Fragment:       f,
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

// BuildArg returns the rendered the bindvar at index.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *FragmentContext) BuildArg(index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if index > len(c.Fragment.Args) {
		return "", fmt.Errorf("%w: bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Args))
	}
	c.Global.onBindVar(defaultStyle)
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.Global.bindVarStyle == syntax.Question {
		c.Global.argStore = append(c.Global.argStore, c.Fragment.Args[i])
		if c.Global.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(c.Global.argStore))
		}
		c.argsBuilt[i] = built
	}
	return built, nil
}

// BuildColumn returns the rendered column at index.
func (c *FragmentContext) BuildColumn(index int) (string, error) {
	if index > len(c.Fragment.Columns) {
		return "", fmt.Errorf("%w: column index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Columns))
	}
	i := index - 1
	c.columnsUsed[i] = true
	col := c.Fragment.Columns[i]
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

// BuildTable returns the rendered table at index.
func (c *FragmentContext) BuildTable(index int) (string, error) {
	if index > len(c.Fragment.Tables) {
		return "", fmt.Errorf("%w: table index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Tables))
	}
	c.tableUsed[index-1] = true
	return string(c.Fragment.Tables[index-1]), nil
}

// BuildFragment returns the rendered fragment at index.
func (c *FragmentContext) BuildFragment(index int) (string, error) {
	if index > len(c.Fragment.Fragments) {
		return "", fmt.Errorf("%w: fragment index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Fragments))
	}
	i := index - 1
	c.fragmentsUsed[i] = true
	seg := c.Fragment.Fragments[i]
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

// BuildBuilder returns the rendered builder at index.
func (c *FragmentContext) BuildBuilder(index int) (string, error) {
	if index > len(c.Fragment.Builders) {
		return "", fmt.Errorf("%w: builder index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Builders))
	}
	i := index - 1
	c.builderUsed[i] = true
	builder := c.Fragment.Builders[i]
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

// ReportUsedArg reports the arg at index is used, starting from 1.
func (c *FragmentContext) ReportUsedArg(index int) {
	if index <= 0 || index > len(c.argsUsed) {
		return
	}
	c.argsUsed[index-1] = true
}

// ReportUsedColumn reports the column at index is used, starting from 1.
func (c *FragmentContext) ReportUsedColumn(index int) {
	if index <= 0 || index > len(c.columnsUsed) {
		return
	}
	c.columnsUsed[index-1] = true
}

// ReportUsedTable reports the table at index is used, starting from 1.
func (c *FragmentContext) ReportUsedTable(index int) {
	if index <= 0 || index > len(c.tableUsed) {
		return
	}
	c.tableUsed[index-1] = true
}

// ReportUsedFragment reports the fragment at index is used, starting from 1.
func (c *FragmentContext) ReportUsedFragment(index int) {
	if index <= 0 || index > len(c.fragmentsUsed) {
		return
	}
	c.fragmentsUsed[index-1] = true
}

// ReportUsedBuilder reports the builder at index is used, starting from 1.
func (c *FragmentContext) ReportUsedBuilder(index int) {
	if index > len(c.builderUsed) {
		return
	}
	c.builderUsed[index-1] = true
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
