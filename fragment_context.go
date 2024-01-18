package sqlf

import (
	"fmt"
	"strings"

	"github.com/go-oauth2/oauth2/v4/errors"
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

// BuildArg returns the built the bindvar at index.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *FragmentContext) BuildArg(index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if index < 1 || index > len(c.Fragment.Args) {
		return "", fmt.Errorf("%w: bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Args))
	}
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.Global.bindVarStyle == syntax.Question {
		built = c.Global.CommitBuiltArg(c.Fragment.Args[i], defaultStyle)
		c.argsBuilt[i] = built
	}
	return built, nil
}

// BuildColumn returns the built column at index.
func (c *FragmentContext) BuildColumn(index int) (string, error) {
	if index < 1 || index > len(c.Fragment.Columns) {
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
	fragment := &Fragment{
		Raw:    column.Raw,
		Args:   column.Args,
		Tables: []Table{column.Table},
	}
	ctxColumn := newFragmentContext(c.Global, fragment)
	built, err := build(ctxColumn)
	if err != nil {
		return "", err
	}
	// don't check usage of tables
	for i := range ctxColumn.tableUsed {
		ctxColumn.tableUsed[i] = true
	}
	if err := ctxColumn.CheckUsage(); err != nil {
		return "", fmt.Errorf("build '%s': %w", column.Raw, err)
	}
	return built, err
}

// BuildTable returns the built table at index.
func (c *FragmentContext) BuildTable(index int) (string, error) {
	if index < 1 || index > len(c.Fragment.Tables) {
		return "", fmt.Errorf("%w: table index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Tables))
	}
	c.tableUsed[index-1] = true
	return string(c.Fragment.Tables[index-1]), nil
}

// BuildFragment returns the built fragment at index.
func (c *FragmentContext) BuildFragment(index int) (string, error) {
	if index < 1 || index > len(c.Fragment.Fragments) {
		return "", fmt.Errorf("%w: fragment index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.Fragment.Fragments))
	}
	i := index - 1
	c.fragmentsUsed[i] = true
	fragment := c.Fragment.Fragments[i]
	built := c.fragmentsBuilt[i]
	if built == "" || (c.Global.bindVarStyle == syntax.Question && len(fragment.Args) > 0) {
		b, err := fragment.BuildContext(c.Global)
		if err != nil {
			return "", err
		}
		c.fragmentsBuilt[i] = b
		built = b
	}
	return built, nil
}

// BuildBuilder returns the built builder at index.
func (c *FragmentContext) BuildBuilder(index int) (string, error) {
	if index < 1 || index > len(c.Fragment.Builders) {
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
	if index < 1 || index > len(c.argsUsed) {
		return
	}
	c.argsUsed[index-1] = true
}

// ReportUsedColumn reports the column at index is used, starting from 1.
func (c *FragmentContext) ReportUsedColumn(index int) {
	if index < 1 || index > len(c.columnsUsed) {
		return
	}
	c.columnsUsed[index-1] = true
}

// ReportUsedTable reports the table at index is used, starting from 1.
func (c *FragmentContext) ReportUsedTable(index int) {
	if index < 1 || index > len(c.tableUsed) {
		return
	}
	c.tableUsed[index-1] = true
}

// ReportUsedFragment reports the fragment at index is used, starting from 1.
func (c *FragmentContext) ReportUsedFragment(index int) {
	if index < 1 || index > len(c.fragmentsUsed) {
		return
	}
	c.fragmentsUsed[index-1] = true
}

// ReportUsedBuilder reports the builder at index is used, starting from 1.
func (c *FragmentContext) ReportUsedBuilder(index int) {
	if index < 1 || index > len(c.builderUsed) {
		return
	}
	c.builderUsed[index-1] = true
}

// CheckUsage checks if all args, columns, tables, fragments and builders are used.
func (c *FragmentContext) CheckUsage() error {
	if c == nil {
		return nil
	}
	msgs := make([]string, 0, 5)
	if msg := unusedIndexs(c.argsUsed); msg != "" {
		msgs = append(msgs, "unused args: ["+msg+"]")
	}
	if msg := unusedIndexs(c.columnsUsed); msg != "" {
		msgs = append(msgs, "unused columns: ["+msg+"]")
	}
	if msg := unusedIndexs(c.tableUsed); msg != "" {
		msgs = append(msgs, "unused tables: ["+msg+"]")
	}
	if msg := unusedIndexs(c.fragmentsUsed); msg != "" {
		msgs = append(msgs, "unused fragments: ["+msg+"]")
	}
	if msg := unusedIndexs(c.builderUsed); msg != "" {
		msgs = append(msgs, "unused builders: ["+msg+"]")
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}

func unusedIndexs(used []bool) string {
	unused := new(strings.Builder)
	for i, v := range used {
		if !v {
			if unused.Len() > 0 {
				unused.WriteString(", ")
			}
			unused.WriteString(fmt.Sprint(i + 1))
		}
	}
	return unused.String()
}
