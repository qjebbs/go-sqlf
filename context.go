package sqls

import (
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqls/syntax"
)

// Context is the global context shared between all segments building.
type Context struct {
	// args store
	ArgStore *[]any
	// override  bindvar style of all segments
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
		return "", fmt.Errorf("invalid bindvar index %d", index)
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

// context is the context for current segment building.
type context struct {
	global  *Context // global context
	Segment *Segment // current segment

	ArgsBuilt     []string // cache of built args
	ColumnsBuilt  []string // cache of built columns
	SegmentsBuilt []string // cache of built segments
	BuildersBuilt []string // cache of built builders

	ArgsUsed     []bool // flags to indicate if an arg is used
	ColumnsUsed  []bool // flags to indicate if a column is used
	TableUsed    []bool // flag to indicate if a table is used
	SegmentsUsed []bool // flags to indicate if a segment is used
	BuilderUsed  []bool // flags to indicate if a builder is used
}

func newSegmentContext(ctx *Context, s *Segment) *context {
	if s == nil {
		return nil
	}
	return &context{
		global:        ctx,
		Segment:       s,
		ArgsBuilt:     make([]string, len(s.Args)),
		ColumnsBuilt:  make([]string, len(s.Columns)),
		TableUsed:     make([]bool, len(s.Tables)),
		SegmentsBuilt: make([]string, len(s.Segments)),
		BuildersBuilt: make([]string, len(s.Builders)),
		ArgsUsed:      make([]bool, len(s.Args)),
		ColumnsUsed:   make([]bool, len(s.Columns)),
		SegmentsUsed:  make([]bool, len(s.Segments)),
		BuilderUsed:   make([]bool, len(s.Builders)),
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
	for i, v := range c.SegmentsUsed {
		if !v {
			return fmt.Errorf("segment %d is not used", i+1)
		}
	}
	for i, v := range c.BuilderUsed {
		if !v {
			return fmt.Errorf("builder %d is not used", i+1)
		}
	}
	return nil
}
