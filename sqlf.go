// Package sqlf is dedicated to building SQL queries by composing fragments.
//
// Unlike other SQL builders or ORMs, *Fragment is the only concept you need
// to understand.
// It uses the same bind variable syntax (? / $1) as database/sql, and also
// supports binding other fragment builders for flexible query composition.
package sqlf

import (
	"context"
)

// Builder is a SQL fragment builder.
type Builder interface {
	// Build builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	Build(ctx *Context) (query string, err error)
}

// BuildQuery builds the given builder into a query string and args slice.
// The default dialect and argument store will be used if not set in the context.
// To customize the dialect or argument store, use ContextWithDialect or ContextWithArgStore
// to create a new context and pass it to this function.
// E.g.:
//
//	ctx = sqlf.ContextWithArgStore(ctx, argstore.NewPositional())
//	ctx = sqlf.ContextWithDialect(ctx, dialect.PostgreSQL{})
//	query, args, err := sqlf.BuildQuery(ctx, builder)
func BuildQuery(ctx context.Context, b Builder) (query string, args []any, err error) {
	if b == nil {
		return "", nil, nil
	}
	// new context with default dialect and argument store if not set
	buildCtx := NewContext(ctx)
	query, err = b.Build(buildCtx)
	if err != nil {
		return "", nil, err
	}
	return query, buildCtx.Args(), nil
}
