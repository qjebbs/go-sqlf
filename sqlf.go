// Package sqlf is dedicated to building SQL queries by composing fragments.
//
// Unlike other SQL builders or ORMs, *Fragment is the only concept you need
// to understand.
// It uses the same bind variable syntax (? / $1) as database/sql, and also
// supports binding other fragment builders for flexible query composition.
package sqlf

// Builder is a SQL fragment builder.
type Builder interface {
	// Build builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	Build(ctx *Context) (query string, err error)
}
