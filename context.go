package sqlf

import (
	"context"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/arg"
)

var _ context.Context = (*Context)(nil)

var defaultDialect dialect.Dialect = dialect.PostgreSQL{}

// Context is the context for fragment building.
type Context struct {
	context.Context
	key, value any
}

// NewContext returns a new Context with an argument store for the given dialect.
// If no store is provided, a new one is created using the dialect's NewArgStore method.
func NewContext(parent context.Context, dialect dialect.Dialect) *Context {
	if dialect == nil {
		dialect = defaultDialect
	}
	ctx := contextWithValue(parent, dialectKey{}, dialect)
	ctx = contextWithValue(ctx, argStoreKey{}, arg.NewArgStoreFromStyle(dialect.BindVarStyle()))
	return ctx
}

// ContextWithValue returns a new context with the given key and value.
// It's used to create from *Context to a new *Context with custom values.
//
// If you want to create a new context.Context with custom values, turn to context.WithValue.
func ContextWithValue(parent *Context, key, value any) *Context {
	// the parent is must of type *Context, so we can avoid dialect and arg checking,
	// since any other path to create a *Context has already ensured those values are set.
	return contextWithValue(parent, key, value)
}

// contextWithValue returns a new context with the given key and value.
func contextWithValue(parent context.Context, key, value any) *Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return &Context{
		Context: parent,
		key:     key,
		value:   value,
	}
}

// Value retrieves the value in the context.
func (c *Context) Value(key any) any {
	if c.key == key {
		return c.value
	}
	return c.Context.Value(key)
}

type argStoreKey struct{}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	store := c.Value(argStoreKey{}).(arg.Store)
	return store.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
func (c *Context) CommitArg(v any) string {
	store := c.Value(argStoreKey{}).(arg.Store)
	return store.CommitArg(v)
}

// ContextWithNewArgStore returns a new context with a new ArgStore created from the dialect in the parent context.
//
// It's useful for creating sub-contexts that need their own ArgStore, like what sqlf.Build() does.
func ContextWithNewArgStore(parent *Context) *Context {
	dialect := parent.Value(dialectKey{}).(dialect.Dialect)
	store := arg.NewArgStoreFromStyle(dialect.BindVarStyle())
	return contextWithArgStore(parent, store)
}

// contextWithArgStore returns a new context with the given ArgStore.
// It panics if the store is nil.
func contextWithArgStore(parent context.Context, store arg.Store) *Context {
	ctx := contextWithValue(parent, argStoreKey{}, store)
	return ctx
}

type dialectKey struct{}

// ContextWithDialect returns a new context with the given dialect.
func ContextWithDialect(parent context.Context, dialect dialect.Dialect) *Context {
	ctx := contextWithValue(parent, dialectKey{}, dialect)
	if ctx.Value(argStoreKey{}) == nil {
		return contextWithArgStore(ctx, arg.NewArgStoreFromStyle(dialect.BindVarStyle()))
	}
	return ctx
}

// DialectFromContext retrieves the dialect from the context.
func DialectFromContext(ctx context.Context) (dialect.Dialect, bool) {
	if d, ok := ctx.Value(dialectKey{}).(dialect.Dialect); ok {
		return d, true
	}
	return nil, false
}

// Dialect returns the dialect of the context.
func (c *Context) Dialect() dialect.Dialect {
	r := c.Value(dialectKey{}).(dialect.Dialect)
	return r
}
