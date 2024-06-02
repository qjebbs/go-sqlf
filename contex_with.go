package sqlf

func contextWith(c *Context, fn func(c *Context) error) (*Context, error) {
	ctx := &Context{
		parent: c,
	}
	return ctx, fn(ctx)
}

func contextValue[T any](c *Context, fn func(c *Context) (T, bool)) (T, bool) {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if v, ok := fn(ctx); ok {
			return v, true
		}
	}
	var v T
	return v, false
}

// root returns the root context.
func (c *Context) root() *Context {
	root := c
	for root.parent != nil {
		root = root.parent
	}
	return root
}
