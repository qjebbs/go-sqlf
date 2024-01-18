package sqlf

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/qjebbs/go-sqlf/syntax"
)

// Context is the global context shared between all fragments building.
type Context struct {
	funcs map[string]*funcInfo

	argStore     []any
	bindVarStyle syntax.BindVarStyle

	globalArgs      []any    // args to be referenced by other builders, with Context.Arg(int)
	globalArgsUsed  []bool   // flags to indicate if an arg is used
	globalArgsBuilt []string // cache of built args
}

// NewContext returns a new context.
func NewContext() *Context {
	funcs, err := createValueFuncs(builtInFuncs)
	if err != nil {
		// should never happen for builtInFuncs
		panic(err)
	}
	return &Context{
		funcs:    funcs,
		argStore: make([]any, 0),
	}
}

// Funcs adds the preprocessing functions of the FuncMap to the context.
//
// The function name is case sensitive, only letters and underscore are allowed.
//
// Allowed function signatures are:
//
//	func(/* args... */)
//	func(/* args... */) string
//	func(/* args... */) (string, error)
//
// Allowed argument types are:
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

// WithBindVarStyle set the bindvar style to the context, which
// overrides bindvar style of all fragments.
// if not, the first bindvar style encountered when building is applied.
func (c *Context) WithBindVarStyle(style syntax.BindVarStyle) *Context {
	c.bindVarStyle = style
	return c
}

// WithGlobalArgs set the global args to the context, which can be referenced by #globalArgDollar and #globalArgQuestion.
func (c *Context) WithGlobalArgs(args []any) *Context {
	c.globalArgs = args
	c.globalArgsBuilt = make([]string, len(args))
	c.globalArgsUsed = make([]bool, len(args))
	return c
}

// BuiltArgs returns the built args of the context.
func (c *Context) BuiltArgs() ([]any, error) {
	if err := c.CheckUsage(); err != nil {
		return nil, err
	}
	return c.argStore, nil
}

// CommitBuiltArg commits the arg to the built args of the context and returns the built bindvar.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *Context) CommitBuiltArg(arg any, defaultStyle syntax.BindVarStyle) string {
	if c.bindVarStyle == 0 {
		c.bindVarStyle = defaultStyle
	}
	c.argStore = append(c.argStore, arg)
	if c.bindVarStyle == syntax.Question {
		return "?"
	}
	return "$" + strconv.Itoa(len(c.argStore))
}

// BuildGlobalArg returns the built globalArg in the context at index.
// defaultStyle is used only when no style is set in the context and no style is seen before.
func (c *Context) BuildGlobalArg(index int, defaultStyle syntax.BindVarStyle) (string, error) {
	if index < 1 || index > len(c.globalArgs) {
		return "", fmt.Errorf("%w: global bindvar index %d out of range [1,%d]", ErrInvalidIndex, index, len(c.globalArgs))
	}
	i := index - 1
	c.globalArgsUsed[i] = true
	built := c.globalArgsBuilt[i]
	if built == "" || c.bindVarStyle == syntax.Question {
		built = c.CommitBuiltArg(c.globalArgs[i], defaultStyle)
		c.globalArgsBuilt[i] = built
	}
	return built, nil
}

// ReportUsedGlobalArg reports the global arg at index is used.
func (c *Context) ReportUsedGlobalArg(index int) {
	if index < 1 || index > len(c.globalArgsUsed) {
		return
	}
	c.globalArgsUsed[index-1] = true
}

// CheckUsage checks if all args are used.
func (c *Context) CheckUsage() error {
	if msg := unusedIndexs(c.globalArgsUsed); msg != "" {
		return errors.New("unused global args: [" + msg + "]")
	}
	return nil
}
