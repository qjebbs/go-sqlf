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

	// cached values
	d dialect.Dialect
	s arg.Store
}

// NewContext returns a new Context with both the dialect and arg store set.
func NewContext(parent context.Context, dialect dialect.Dialect) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if dialect == nil {
		dialect = defaultDialect
	}
	store := arg.NewArgStoreFromStyle(dialect.BindVarStyle())
	ctx := context.WithValue(parent, dialectKey{}, dialect)
	ctx = context.WithValue(ctx, argStoreKey{}, store)
	return &defaultCtx{
		Context: ctx,
		d:       dialect,
		s:       store,
	}
}

// ContextWithValue returns a new Context derived from parent with the key and value set.
//
// Use this when creating a Context from an existing Context.
// If you want to create from context.Context to context.Context, use context.WithValue.
func ContextWithValue(parent Context, key, value any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	// the parent must of type Context, so we can avoid dialect and arg checking,
	// since any other path to create a Context has already ensured those values are set.
	return &defaultCtx{
		Context: context.WithValue(unwrapContext(parent), key, value),
	}
}

// unwrapContext extracts *defaultCtx.Context to avoid double wrapping.
func unwrapContext(ctx Context) context.Context {
	if ctx, ok := ctx.(*defaultCtx); ok {
		// optimize for the common case
		return ctx.Context
	}
	return ctx
}

// ContextWithNewArgStore returns a new context with a new ArgStore created from the dialect in the parent context.
//
// It's useful for creating sub-contexts that need their own ArgStore, like what sqlf.Build() does.
func ContextWithNewArgStore(parent Context) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	store := arg.NewArgStoreFromStyle(parent.BaseDialect().BindVarStyle())
	return ContextWithValue(parent, argStoreKey{}, store)
}

// BaseDialect implemens the Context interface.
func (c *defaultCtx) BaseDialect() dialect.Dialect {
	// no need to check nil c, since user cannot create defaultCtx directly.
	if c.d != nil {
		return c.d
	}
	// no need to check existence, since NewContext always sets it.
	c.d = c.Value(dialectKey{}).(dialect.Dialect)
	return c.d
}

// Args implements the Context interface.
func (c *defaultCtx) Args() []any {
	return c.store().Args()
}

// CommitArg implements the Context interface.
func (c *defaultCtx) CommitArg(v any) string {
	return c.store().CommitArg(v)
}

func (c *defaultCtx) store() arg.Store {
	if c.s != nil {
		return c.s
	}
	c.s = c.Value(argStoreKey{}).(arg.Store)
	return c.s
}
