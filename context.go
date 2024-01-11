package sqlf

import (
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	// args store
	ArgStore *[]any
	// override  bindvar style of all fragments
	BindVarStyle syntax.BindVarStyle

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

// WithArgs set the args to the context.
func (c *Context) WithArgs(args []any) *Context {
	c.args = args
	c.argsBuilt = make([]string, len(args))
	c.argsUsed = make([]bool, len(args))
	return c
}

// Arg returns the built arg in the context at index.
func (c *Context) Arg(index int) (string, error) {
	if index > len(c.args) {
		return "", fmt.Errorf("%w: global bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.args))
	}
	if c.BindVarStyle == 0 {
		c.BindVarStyle = syntax.Dollar
	}
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.BindVarStyle == syntax.Question {
		*c.ArgStore = append(*c.ArgStore, c.args[i])
		if c.BindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(*c.ArgStore))
		}
		c.argsBuilt[i] = built
	}
	return built, nil
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

func newFragmentContext(ctx *Context, s *Fragment) *context {
	if s == nil {
		return nil
	}
	return &context{
		global:         ctx,
		Fragment:       s,
		ArgsBuilt:      make([]string, len(s.Args)),
		ColumnsBuilt:   make([]string, len(s.Columns)),
		TableUsed:      make([]bool, len(s.Tables)),
		FragmentsBuilt: make([]string, len(s.Fragments)),
		BuildersBuilt:  make([]string, len(s.Builders)),
		ArgsUsed:       make([]bool, len(s.Args)),
		ColumnsUsed:    make([]bool, len(s.Columns)),
		FragmentsUsed:  make([]bool, len(s.Fragments)),
		BuilderUsed:    make([]bool, len(s.Builders)),
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
