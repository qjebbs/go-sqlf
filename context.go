package sqlf

import (
	"context"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/arg"
)

var _ Context = (*defaultCtx)(nil)

var defaultDialect dialect.Dialect = dialect.PostgreSQL{}

type argStoreKey struct{}
type dialectKey struct{}

type defaultCtx struct {
	context.Context
}

// NewContext returns a new Context with both the dialect and arg store set.
func NewContext(parent context.Context, dialect dialect.Dialect) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if dialect == nil {
		dialect = defaultDialect
	}
	ctx := context.WithValue(parent, dialectKey{}, dialect)
	ctx = context.WithValue(ctx, argStoreKey{}, arg.NewArgStoreFromStyle(dialect.BindVarStyle()))
	return &defaultCtx{ctx}
}

// ContextWithValue returns a new Context derived from parent with the key and value set.
//
// Use this when creating a Context from an existing Context.
// If you want to create from context.Context to context.Context, use context.WithValue.
func ContextWithValue(parent Context, key, value any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if ctx, ok := parent.(*defaultCtx); ok {
		// optimize for the common case
		return &defaultCtx{
			context.WithValue(ctx.Context, key, value),
		}
	}
	// the parent must of type Context, so we can avoid dialect and arg checking,
	// since any other path to create a Context has already ensured those values are set.
	return &defaultCtx{
		context.WithValue(parent, key, value),
	}
}

// ContextWithNewArgStore returns a new context with a new ArgStore created from the dialect in the parent context.
//
// It's useful for creating sub-contexts that need their own ArgStore, like what sqlf.Build() does.
func ContextWithNewArgStore(parent Context) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	store := arg.NewArgStoreFromStyle(parent.Dialect().BindVarStyle())
	return ContextWithValue(parent, argStoreKey{}, store)
}

// Dialect returns the dialect of the context.
func (c *defaultCtx) Dialect() dialect.Dialect {
	// no need to check nil, since user cannot create _Context directly.
	// no need to check existence, since NewContext always sets it.
	return c.Value(dialectKey{}).(dialect.Dialect)
}

// Args returns the built args of the context.
func (c *defaultCtx) Args() []any {
	store := c.Value(argStoreKey{}).(arg.Store)
	return store.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
func (c *defaultCtx) CommitArg(v any) string {
	store := c.Value(argStoreKey{}).(arg.Store)
	return store.CommitArg(v)
}
