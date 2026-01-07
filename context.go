package sqlf

import (
	"context"
	"time"
)

var _ context.Context = (*Context)(nil)

// Context is the context for fragment building.
type Context struct {
	parent     context.Context
	key, value any
}

// NewContext returns a new Context with an argument store for the given bind style.
func NewContext(parent context.Context, style BindStyle) *Context {
	return &Context{
		parent: parent,
		key:    argStoreKey{},
		value:  newArgStore(style),
	}
}

// Deadline always returns false.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.parent != nil {
		return c.parent.Deadline()
	}
	return
}

// Done always returns nil.
func (c *Context) Done() <-chan struct{} {
	if c.parent != nil {
		return c.parent.Done()
	}
	return nil
}

// Err always returns nil.
func (c *Context) Err() error {
	if c.parent != nil {
		return c.parent.Err()
	}
	return nil
}

// Value retrieves the value in the context.
func (c *Context) Value(key any) any {
	if c.key == key {
		return c.value
	}
	if c.parent != nil {
		return c.parent.Value(key)
	}
	return nil
}
