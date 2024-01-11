package sqlf

import (
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	// ArgStore is the storage for built args.
	ArgStore *[]any

	bindVarStyle syntax.BindVarStyle

	args      []any    // args to be referenced by other builders, with Context.Arg(int)
	argsUsed  []bool   // flags to indicate if an arg is used
	argsBuilt []string // cache of built args
}

// NewContext returns a new context.
func NewContext(argStore *[]any) *Context {
	return &Context{
		ArgStore: argStore,
	}
}

// WithArgs set the args to the context, which can be referenced by #globalArgDollar and #globalArgQuestion.
//
// Note: it has nothing to do with the c.ArgStore.
func (c *Context) WithArgs(args []any) *Context {
	c.args = args
	c.argsBuilt = make([]string, len(args))
	c.argsUsed = make([]bool, len(args))
	return c
}

// WithBindVarStyle set the bindvar style to the context, which
// overrides bindvar style of all fragments.
// if not, the first bindvar style encountered when building is applied.
func (c *Context) WithBindVarStyle(style syntax.BindVarStyle) *Context {
	c.bindVarStyle = style
	return c
}

// buildArg returns the built buildArg in the context at index.
func (c *Context) buildArg(index int, style syntax.BindVarStyle) (string, error) {
	if index > len(c.args) {
		return "", fmt.Errorf("%w: global bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.args))
	}
	c.onBindArg(style)
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.bindVarStyle == syntax.Question {
		*c.ArgStore = append(*c.ArgStore, c.args[i])
		if c.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(*c.ArgStore))
		}
		c.argsBuilt[i] = built
	}
	return built, nil
}

func (c *Context) onBindArg(style syntax.BindVarStyle) {
	if c.bindVarStyle == 0 {
		c.bindVarStyle = style
	}
}

func (c *Context) checkUsage() error {
	for i, v := range c.argsUsed {
		if !v {
			return fmt.Errorf("arg %d is not used", i+1)
		}
	}
	return nil
}

// context is the context for current fragment building.
type context struct {
	global   *Context  // global context
	Fragment *Fragment // current fragment

	ArgsBuilt      []string // cache of built args
	ColumnsBuilt   []string // cache of built columns
	FragmentsBuilt []string // cache of built fragments
	BuildersBuilt  []string // cache of built builders

	ArgsUsed      []bool // flags to indicate if an arg is used
	ColumnsUsed   []bool // flags to indicate if a column is used
	TableUsed     []bool // flag to indicate if a table is used
	FragmentsUsed []bool // flags to indicate if a fragment is used
	BuilderUsed   []bool // flags to indicate if a builder is used
}

func newFragmentContext(ctx *Context, f *Fragment) *context {
	if f == nil {
		return nil
	}
	return &context{
		global:         ctx,
		Fragment:       f,
		ArgsBuilt:      make([]string, len(f.Args)),
		ColumnsBuilt:   make([]string, len(f.Columns)),
		TableUsed:      make([]bool, len(f.Tables)),
		FragmentsBuilt: make([]string, len(f.Fragments)),
		BuildersBuilt:  make([]string, len(f.Builders)),
		ArgsUsed:       make([]bool, len(f.Args)),
		ColumnsUsed:    make([]bool, len(f.Columns)),
		FragmentsUsed:  make([]bool, len(f.Fragments)),
		BuilderUsed:    make([]bool, len(f.Builders)),
	}
}

func (c *context) checkUsage() error {
	if c == nil {
		return nil
	}
	for i, v := range c.ArgsUsed {
		if !v {
			return fmt.Errorf("arg %d is not used", i+1)
		}
	}
	for i, v := range c.ColumnsUsed {
		if !v {
			return fmt.Errorf("column %d is not used", i+1)
		}
	}
	for i, v := range c.TableUsed {
		if !v {
			return fmt.Errorf("table %d is not used", i+1)
		}
	}
	for i, v := range c.FragmentsUsed {
		if !v {
			return fmt.Errorf("fragment %d is not used", i+1)
		}
	}
	for i, v := range c.BuilderUsed {
		if !v {
			return fmt.Errorf("builder %d is not used", i+1)
		}
	}
	return nil
}
