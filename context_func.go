package sqlf

// Funcs adds the preprocessing functions to the context.
// func (c *Context) Funcs(funcs FuncMap) error {
// 	return addValueFuncs(c.funcs, funcs)
// }

// ContextWithFuncs returns a new context with the preprocessing functions added.
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
//   - *sqlf.Context: allowed only as the first argument
//
// Here are examples of legal names and function signatures:
//
//	funcs := sqlf.FuncMap{
//		// #number1, #join('#number', ', ')
//		"number": func(i int) (string, error) {/* ... */},
//		// #myBuilder1, #join('#myBuilder', ', ')
//		"myBuilder": func(ctx *sqlf.Context, i int) (string, error)  {/* ... */},
//		// #string('string')
//		"string": func(str string) (string, error)  {/* ... */},
//		// #numbers(1,2)
//		"numbers": func(ctx *sqlf.Context, a, b int) string  {/* ... */},
//	}
func ContextWithFuncs(c *Context, funcs FuncMap) (*Context, error) {
	ctx := &Context{
		parent: c,
		funcs:  make(map[string]*funcInfo),
	}
	err := addValueFuncs(ctx.funcs, funcs)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (c *Context) fn(name string) (*funcInfo, bool) {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if ctx.funcs == nil {
			continue
		}
		if fn, ok := ctx.funcs[name]; ok {
			return fn, true
		}
	}
	return nil, false
}
