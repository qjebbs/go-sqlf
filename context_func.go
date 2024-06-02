package sqlf

// ContextWithFuncs returns a new context with the preprocessing functions added.
func ContextWithFuncs(c *Context, funcs FuncMap) (*Context, error) {
	return contextWith(c, func(c *Context) error {
		c.funcs = make(map[string]*funcInfo)
		return addValueFuncs(c.funcs, funcs)
	})
}

func (c *Context) fn(name string) (*funcInfo, bool) {
	return contextValue(c, func(c *Context) (*funcInfo, bool) {
		if c.funcs == nil {
			return nil, false
		}
		fn, ok := c.funcs[name]
		return fn, ok
	})
}
