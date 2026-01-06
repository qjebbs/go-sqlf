package sqlf

// Context is the context for fragment building.
type Context struct {
	parent     *Context
	key, value any
}

// NewContext returns a new context.
func NewContext(style BindStyle) *Context {
	return &Context{
		key:   argStoreKey{},
		value: newArgStore(style),
	}
}

// ContextWith returns a new context with the given key and value.
func ContextWith(ctx *Context, key, value any) *Context {
	return &Context{
		parent: ctx,
		key:    key,
		value:  value,
	}
}

// Value retrieves the value in the context.
func (c *Context) Value(key any) any {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if ctx.key == key {
			return ctx.value
		}
	}
	return nil
}
