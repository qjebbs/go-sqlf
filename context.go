package sqlf

import (
	"context"
	"time"

	"github.com/qjebbs/go-sqlf/v4/argstore"
	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ context.Context = (*Context)(nil)

var defaultDialect = dialect.AnsiSQL{}

// Context is the context for fragment building.
type Context struct {
	parent     context.Context
	key, value any
}

// NewContext returns a new Context with an argument store for the given bind style.
func NewContext(parent context.Context) *Context {
	return EnsureContextValues(parent, defaultDialect)
}

// ContextWithValue returns a new context with the given key and value.
func ContextWithValue(parent context.Context, key, value any) *Context {
	return EnsureContextValues(&Context{
		parent: parent,
		key:    key,
		value:  value,
	}, defaultDialect)
}

// EnsureContextValues ensures that the context has both a dialect and an ArgStore.
func EnsureContextValues(parent context.Context, defaultDialect dialect.Dialect) *Context {
	dialect, ctx := contextWithDefaultValue(parent, dialectKey{}, defaultDialect)
	_, ctx = contextWithDefaultValue(ctx, argStoreKey{}, dialect.NewArgStore())
	return ctx
}

// contextWithDefaultValue returns a new context with the given key and value
// only if the key is not already set in the context.
func contextWithDefaultValue[T any](ctx context.Context, key any, value T) (applied T, c *Context) {
	if v := ctx.Value(key); v != nil {
		if existingCtx, ok := ctx.(*Context); ok {
			return v.(T), existingCtx
		}
		return v.(T), &Context{
			parent: ctx,
		}
	}
	return value, &Context{
		parent: ctx,
		key:    key,
		value:  value,
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

type argStoreKey struct{}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	store := c.Value(argStoreKey{}).(argstore.Store)
	return store.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
func (c *Context) CommitArg(arg any) string {
	store := c.Value(argStoreKey{}).(argstore.Store)
	return store.CommitArg(arg)
}

// ContextWithArgStore returns a new context with the given ArgStore.
func ContextWithArgStore(ctx context.Context, store argstore.Store) *Context {
	if store == nil {
		panic("store cannot be nil")
	}
	return ContextWithValue(ctx, argStoreKey{}, store)
}

type dialectKey struct{}

// ContextWithDialect returns a new context with the given dialect.
func ContextWithDialect(ctx context.Context, dialect dialect.Dialect) *Context {
	return ContextWithValue(ctx, dialectKey{}, dialect)
}

// DialectFromContext retrieves the dialect from the context.
func DialectFromContext(ctx context.Context) (dialect.Dialect, bool) {
	value := ctx.Value(dialectKey{})
	if value == nil {
		return nil, false
	}
	if d, ok := value.(dialect.Dialect); ok {
		return d, true
	}
	return nil, false
}

// Dialect returns the dialect of the context.
func (c *Context) Dialect() dialect.Dialect {
	r := c.Value(dialectKey{}).(dialect.Dialect)
	return r
}
