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
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *Context) CommitArg(arg any) string {
	root := c.root()
	*root.argStore = append(*root.argStore, arg)
	if root.bindVarStyle == syntax.Question {
		return "?"
	}
	return "$" + strconv.Itoa(len(*root.argStore))
}
