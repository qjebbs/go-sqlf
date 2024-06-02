package sqlf

import (
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
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
