// Package sqlf focuses only on building SQL queries by combining fragments.
//
// The package exports only a few functions and methods, but improves a lot on the
// reusability and extensibility of SQL, which are the main challenges we encounter
// when writing SQL.
//
// # Fragment
//
// Unlike any other sql builder or ORMs, *Fragment is the only concept you need to learn.
//
// A *Fragment (created by F()) is usually part of a SQL query, which has exactly the same bind
// var syntax (? / $1) as database/sql, but more than that, it allows you
// to bind other fragment builders.
package sqlf

// Builder is a SQL fragment builder.
type Builder interface {
	// Build builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	Build(ctx *Context) (query string, err error)
}
