package sqlf

// Value returns the value in the context.
func (c *Context) Value(key any) any {
	return c.root().values[key]
}

// WithValue sets the value in the context.
func (c *Context) WithValue(key, value any) *Context {
	c.root().values[key] = value
	return c
}
