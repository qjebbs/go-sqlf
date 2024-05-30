package sqlf

// Funcs adds the preprocessing functions to the context.
// func (c *Context) Funcs(funcs FuncMap) error {
// 	return addValueFuncs(c.funcs, funcs)
// }

// ContextWithFuncs returns a new context with the preprocessing functions added.
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
