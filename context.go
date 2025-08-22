package sqlf

import (
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	parent *Context

	bindVarStyle syntax.BindVarStyle
	argStore     argStore

	frag *fragmentContext
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
		argStore: argStore,
	}
}
