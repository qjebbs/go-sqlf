package sqlf

import (
	"context"
	"fmt"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/arg"
)

var defaultDialect dialect.Dialect = dialect.PostgreSQL{}

// NewContext returns a new Context with both the dialect and arg store set.
func NewContext(parent context.Context, dialect dialect.Dialect) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if dialect == nil {
		dialect = defaultDialect
	}
	return newDeafultCtx(parent, dialect)
}

type contextConstraint interface {
	Context
	comparable
}

// ContextWithValue returns a new Context derived from parent with the key and value set.
//
// The new context does not dowgrade the parent context.
//
//	var ctx ExtendedContext // implements sqlf.Context
//	ctx = NewExtendedContext(...)
//	ctx = sqlf.ContextWithValue(ctx, ...) // still a ExtendedContext
func ContextWithValue[T contextConstraint](parent T, key, value any) T {
	var zero T
	if parent == zero {
		panic("cannot create context from nil parent")
	}
	return contextWithValue(parent, key, value)
}

// ContextWithNewArgStore returns a new context with a new ArgStore created from the dialect in the parent context.
//
// The new context does not dowgrade the parent context.
//
//	var ctx ExtendedContext // implements sqlf.Context
//	ctx = NewExtendedContext(...)
//	ctx = sqlf.ContextWithNewArgStore(ctx) // still a ExtendedContext
func ContextWithNewArgStore[T contextConstraint](parent T) T {
	var zero T
	if parent == zero {
		panic("cannot create context from nil parent")
	}
	store := arg.NewArgStoreFromStyle(parent.BaseDialect().BindVarStyle())
	// MUST use parent.ContextWithValue to avoid context downgrading
	return contextWithValue(parent, argStoreKey{}, store)
}

func contextWithValue[T contextConstraint](parent T, key, value any) T {
	// MUST use parent.ContextWithValue to avoid context downgrading
	newCtx := parent.ContextWithValue(key, value)
	ctx, ok := newCtx.(T)
	if !ok {
		panic(fmt.Errorf("%T.ContextWithValue returns %T which does not implement the expected interface", parent, newCtx))
	}
	return ctx
}

var _ Context = (*defaultCtx)(nil)

type defaultCtx struct {
	context.Context

	// cached values
	d dialect.Dialect
	s arg.Store
}

type argStoreKey struct{}
type dialectKey struct{}

func newDeafultCtx(parent context.Context, dialect dialect.Dialect) *defaultCtx {
	store := arg.NewArgStoreFromStyle(dialect.BindVarStyle())
	ctx := context.WithValue(parent, dialectKey{}, dialect)
	ctx = context.WithValue(ctx, argStoreKey{}, store)
	return &defaultCtx{
		Context: ctx,
		d:       dialect,
		s:       store,
	}
}

// ContextWithValue implements the Context interface.
func (c *defaultCtx) ContextWithValue(key, value any) Context {
	return &defaultCtx{
		Context: context.WithValue(c.Context, key, value),
	}
}

// BaseDialect implemens the Context interface.
func (c *defaultCtx) BaseDialect() dialect.Dialect {
	// no need to check nil c, since user cannot create defaultCtx directly.
	if c.d == nil {
		// no need to check existence, since newDeafultCtx always sets it.
		c.d = c.Value(dialectKey{}).(dialect.Dialect)
	}
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
	if c.s == nil {
		c.s = c.Value(argStoreKey{}).(arg.Store)
	}
	return c.s
}
