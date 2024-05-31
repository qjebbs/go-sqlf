package sqlf

import (
	"strconv"

	"github.com/qjebbs/go-sqlf/v2/syntax"
)

// Args returns the built args of the context.
func (c *Context) Args() []any {
	return *c.root().argStore
}

// CommitArg commits an built arg to the context and returns the built bindvar.
//
// It's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) CommitArg(arg any) string {
	root := c.root()
	*root.argStore = append(*root.argStore, arg)
	if root.bindVarStyle == syntax.Question {
		return "?"
	}
	return "$" + strconv.Itoa(len(*root.argStore))
}
