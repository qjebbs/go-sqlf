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

	argStore     []any
	bindVarStyle syntax.BindVarStyle

	globalArgs      []any    // args to be referenced by other builders, with Context.Arg(int)
	globalArgsUsed  []bool   // flags to indicate if an arg is used
	globalArgsBuilt []string // cache of built args
}

// NewContext returns a new context.
func NewContext() *Context {
	return &Context{
		funcs:    createValueFuncs(builtInFuncs),
		argStore: make([]any, 0),
	}
}

// Funcs adds the elements of the argument map to the FuncMap.
func (c *Context) Funcs(funcs FuncMap) *Context {
	addValueFuncs(c.funcs, funcs)
	return c
}

// WithBindVarStyle set the bindvar style to the context, which
// overrides bindvar style of all fragments.
// if not, the first bindvar style encountered when building is applied.
func (c *Context) WithBindVarStyle(style syntax.BindVarStyle) *Context {
	c.bindVarStyle = style
	return c
}

// WithGlobalArgs set the global args to the context, which can be referenced by #globalArgDollar and #globalArgQuestion.
//
// Note: it has nothing to do with the c.ArgStore.
func (c *Context) WithGlobalArgs(args []any) *Context {
	c.globalArgs = args
	c.globalArgsBuilt = make([]string, len(args))
	c.globalArgsUsed = make([]bool, len(args))
	return c
}

// BuiltArgs returns the built args of the context.
func (c *Context) BuiltArgs() []any {
	return c.argStore
}

// BuildArg returns the built BuildArg in the context at index.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *Context) BuildArg(index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if index > len(c.globalArgs) {
		return "", fmt.Errorf("%w: global bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.globalArgs))
	}
	c.onBindVar(defaultStyle)
	i := index - 1
	c.globalArgsUsed[i] = true
	built := c.globalArgsBuilt[i]
	if built == "" || c.bindVarStyle == syntax.Question {
		c.argStore = append(c.argStore, c.globalArgs[i])
		if c.bindVarStyle == syntax.Question {
			built = "?"
		} else {
			built = "$" + strconv.Itoa(len(c.argStore))
		}
		c.globalArgsBuilt[i] = built
	}
	return built, nil
}

func (c *Context) onBindVar(style syntax.BindVarStyle) {
	if c.bindVarStyle == 0 {
		c.bindVarStyle = style
	}
}

// ReportUsedArg reports the global arg at index is used.
func (c *Context) ReportUsedArg(index int) {
	if index > len(c.globalArgsUsed) {
		return
	}
	c.globalArgsUsed[index-1] = true
}

func (c *Context) checkUsage() error {
	for i, v := range c.globalArgsUsed {
		if !v {
			return fmt.Errorf("arg %d is not used", i+1)
		}
	}
	return nil
}
