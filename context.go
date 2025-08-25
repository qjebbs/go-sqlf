package sqlf

import (
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	global *globalContext

	parent     *Context
	key, value any
}

type globalContext struct {
	bindVarStyle syntax.BindVarStyle
	argStore     argStore
}

// NewContext returns a new context.
func NewContext(bindVarStyle syntax.BindVarStyle) *Context {
	ctx := newEmptyContext(bindVarStyle)
	ctx.global.bindVarStyle = bindVarStyle
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
		global: &globalContext{
			bindVarStyle: bindVarStyle,
			argStore:     argStore,
		},
	}
}

// BindVarStyle returns the bind var style of the context.
func (c *Context) BindVarStyle() syntax.BindVarStyle {
	return c.global.bindVarStyle
}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	return c.global.argStore.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
//
// It's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) CommitArg(arg any) string {
	return c.global.argStore.CommitArg(arg)
}

// ContextWith returns a new context with the given key and value.
func ContextWith(ctx *Context, key, value any) *Context {
	newCtx := &Context{
		parent: ctx,
		global: ctx.global,
	}
	newCtx.key = key
	newCtx.value = value
	return newCtx
}

// Value returns the value in the context.
func (c *Context) Value(key any) any {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if ctx.key == key {
			return ctx.value
		}
	}
	return nil
}
