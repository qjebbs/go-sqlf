package sqlf

import (
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	bindVarStyle syntax.BindVarStyle
	argStore     argStore
	values       map[any]any
}

// NewContext returns a new context.
func NewContext(bindVarStyle syntax.BindVarStyle) *Context {
	ctx := newEmptyContext(bindVarStyle)
	ctx.bindVarStyle = bindVarStyle
	return ctx
}

func newEmptyContext(bindVarStyle syntax.BindVarStyle) *Context {
	var argStore argStore
	if bindVarStyle == syntax.Dollar {
		argStore = newDollarArgStore()
	} else {
		argStore = newQuestionArgStore()
	}
	return &Context{
		bindVarStyle: bindVarStyle,
		argStore:     argStore,
		values:       make(map[any]any),
	}
}

// BindVarStyle returns the bind var style of the context.
func (c *Context) BindVarStyle() syntax.BindVarStyle {
	return c.bindVarStyle
}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	return c.argStore.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
//
// It's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) CommitArg(arg any) string {
	return c.argStore.CommitArg(arg)
}

// Value returns the value in the context.
func (c *Context) Value(key any) any {
	return c.values[key]
}

// WithValue sets the value in the context.
func (c *Context) WithValue(key, value any) {
	c.values[key] = value
}
