package sqlf

import (
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	// BindVarStyle overrides bindvar styles of all fragments.
	// if not set, the first bindvar style encountered when
	// building is applied.
	bindVarStyle syntax.BindVarStyle
	argStore     *[]any

	parent   *Context
	funcs    map[string]*funcInfo
	fragment *FragmentContext
}

// NewContext returns a new context.
func NewContext() *Context {
	ctx := newEmptyContext()
	err := addValueFuncs(ctx.funcs, builtInFuncs)
	if err != nil {
		// should never happen for builtInFuncs
		panic(err)
	}
	return ctx
}

// BindVarStyle returns the bindvar style of the context.
func (c *Context) BindVarStyle() syntax.BindVarStyle {
	return c.Root().bindVarStyle
}

// SetBindVarStyle sets the bindvar style of the context.
func (c *Context) SetBindVarStyle(style syntax.BindVarStyle) {
	c.Root().bindVarStyle = style
}

// A

func newEmptyContext() *Context {
	argStore := make([]any, 0)
	return &Context{
		funcs:    make(map[string]*funcInfo),
		argStore: &argStore,
	}
}

// Root returns the root context.
func (c *Context) Root() *Context {
	root := c
	for root.parent != nil {
		root = root.parent
	}
	return root
}

// Parent returns the parent context.
func (c *Context) Parent() *Context {
	return c.parent
}
