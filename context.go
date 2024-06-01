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
	argStore     argStore

	parent *Context
	funcs  map[string]*funcInfo
	frag   *FragmentContext
}

// NewContext returns a new context.
func NewContext(bindVarStyle syntax.BindVarStyle) *Context {
	ctx := newEmptyContext(bindVarStyle)
	ctx.bindVarStyle = bindVarStyle
	err := addValueFuncs(ctx.funcs, builtInFuncs)
	if err != nil {
		// should never happen for builtInFuncs
		panic(err)
	}
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
		funcs:    make(map[string]*funcInfo),
		argStore: argStore,
	}
}

// root returns the root context.
func (c *Context) root() *Context {
	root := c
	for root.parent != nil {
		root = root.parent
	}
	return root
}
