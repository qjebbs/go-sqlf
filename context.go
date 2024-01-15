package sqlf

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	funcs map[string]reflect.Value

	argStore     *[]any
	bindVarStyle syntax.BindVarStyle

	args      []any    // args to be referenced by other builders, with Context.Arg(int)
	argsUsed  []bool   // flags to indicate if an arg is used
	argsBuilt []string // cache of built args
}

// NewContext returns a new context.
func NewContext() *Context {
	argStore := make([]any, 0)
	return &Context{
		funcs:    createValueFuncs(builtInFuncs),
		argStore: &argStore,
	}
}

// BuiltArgs returns the built args of the context.
func (c *Context) BuiltArgs() []any {
	return *c.argStore
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

// Funcs adds the elements of the argument map to the FuncMap.
func (c *Context) Funcs(funcs FuncMap) {
	addValueFuncs(c.funcs, funcs)
}

// Arg returns the built Arg in the context at index.
func (c *Context) Arg(index int, style syntax.BindVarStyle) (string, error) {
	if index > len(c.args) {
		return "", fmt.Errorf("%w: global bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.args))
	}
	c.onBindArg(style)
	i := index - 1
	c.argsUsed[i] = true
	built := c.argsBuilt[i]
	if built == "" || c.bindVarStyle == syntax.Question {
		*c.argStore = append(*c.argStore, c.args[i])
		if c.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(*c.argStore))
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
