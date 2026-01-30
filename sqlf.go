// Package sqlf provides a simple way to build SQL queries by composing fragments.
//
// The core concept is the Fragment, which represents a piece of SQL with
// arguments. Fragments can be composed to build complex queries. It uses the
// same bind variable syntax (? or $N) as database/sql.
package sqlf

import (
	"context"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/arg"
)

// Builder is the interface implemented by types that can build themselves
// into a SQL query string.
type Builder interface {
	// BuildTo builds the SQL fragment for the current builder based on the given
	// context, and commits any arguments to the context.
	BuildTo(ctx Context) (query string, err error)
}

// Context is the context for fragment building.
type Context interface {
	context.Context
	// ContextWithValue returns a NEW Context carrying the given key/value.
	//
	// Implementations MUST NOT mutate the receiver and SHOULD preserve any
	// current-level extensions so the returned Context does not get "downgraded".
	//
	// Callers should normally use sqlf.ContextWithValue to ensure type preservation.
	ContextWithValue(key, value any) Context
	// BaseDialect returns the base dialect of the context,
	// which is the minimal implementation required by sqlf package.
	BaseDialect() dialect.Dialect
	// Args returns the committed arguments in order.
	Args() []any
	// CommitArg commits an argument to the context and returns the built bindvar.
	CommitArg(v any) string
}

// Build builds a Builder into a query string and its corresponding arguments.
// It creates a new context to ensure that the original context is not modified.
func Build[T contextConstraint](ctx T, b Builder) (query string, args []any, err error) {
	if b == nil {
		return "", nil, nil
	}
	// make sure not committing args to the original context
	store := arg.NewArgStoreFromStyle(ctx.BaseDialect().BindVarStyle())
	ctx = ContextWithValue(ctx, argStoreKey{}, store)
	query, err = b.BuildTo(ctx)
	if err != nil {
		return "", nil, err
	}
	return query, ctx.Args(), nil
}
