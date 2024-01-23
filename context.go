package sqlf

import (
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	// BindVarStyle overrides bindvar styles of all fragments.
	// if not set, the first bindvar style encountered when
	// building is applied.
	BindVarStyle syntax.BindVarStyle

	funcs    map[string]*funcInfo
	argStore []any
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

func newEmptyContext() *Context {
	return &Context{
		funcs:    make(map[string]*funcInfo),
		argStore: make([]any, 0),
	}
}

func (c *Context) fn(name string) (*funcInfo, bool) {
	if c == nil || c.funcs == nil {
		return nil, false
	}
	fn, ok := c.funcs[name]
	return fn, ok
}

// Funcs adds the preprocessing functions to the context.
//
// The function name is case sensitive, only letters and underscore are allowed.
//
// Allowed function signatures:
//
//	func(/* args... */) (string, error)
//	func(/* args... */) string
//	func(/* args... */)
//
// Allowed argument types:
//   - number types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,float32, float64
//   - string
//   - bool
//   - *sqlf.FragmentContext: allowed only as the first argument
//
// Here are examples of legal names and function signatures:
//
//	funcs := sqlf.FuncMap{
//		// #number1, #join('#number', ', ')
//		"number": func(i int) (string, error) {/* ... */},
//		// #myBuilder1, #join('#myBuilder', ', ')
//		"myBuilder": func(ctx *sqlf.FragmentContext, i int) (string, error)  {/* ... */},
//		// #string('string')
//		"string": func(str string) (string, error)  {/* ... */},
//		// #numbers(1,2)
//		"numbers": func(ctx *sqlf.FragmentContext, a, b int) string  {/* ... */},
//	}
func (c *Context) Funcs(funcs FuncMap) error {
	return addValueFuncs(c.funcs, funcs)
}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	return c.argStore
}

// CommitArg commits an built arg to the context and returns the built bindvar.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *Context) CommitArg(arg any, defaultStyle syntax.BindVarStyle) string {
	if c.BindVarStyle == 0 {
		c.BindVarStyle = defaultStyle
	}
	c.argStore = append(c.argStore, arg)
	if c.BindVarStyle == syntax.Question {
		return "?"
	}
	return "$" + strconv.Itoa(len(c.argStore))
}
