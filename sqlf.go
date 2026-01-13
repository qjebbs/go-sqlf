// Package sqlf is dedicated to building SQL queries by composing fragments.
//
// Unlike other SQL builders or ORMs, *Fragment is the only concept you need
// to understand.
// It uses the same bind variable syntax (? / $1) as database/sql, and also
// supports binding other fragment builders for flexible query composition.
package sqlf

// Builder is a SQL fragment builder.
type Builder interface {
	// BuildTo builds current fragment into the given context, together with other fragments.
	// The args should be committed to the ctx if any.
	BuildTo(ctx *Context) (query string, err error)
}

// Build builds the given builder into a query string and args slice.
// Unlike Builder.BuildTo, this function does not commit args to the ctx.
func Build(ctx *Context, b Builder) (query string, args []any, err error) {
	if b == nil {
		return "", nil, nil
	}
	// make sure not committing args to the original context
	buildCtx := ContextWithArgStore(ctx, ctx.Dialect().NewArgStore())
	query, err = b.BuildTo(buildCtx)
	if err != nil {
		return "", nil, err
	}
	return query, buildCtx.Args(), nil
}
