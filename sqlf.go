// Package sqlf focuses only on building SQL queries by combining fragments.
//
// The package exports only a few functions and methods, but improves a lot on the
// reusability and extensibility of SQL, which are the main challenges we encounter
// when writing SQL.
//
// # Fragment
//
// Unlike any other sql builder or ORMs, `*Fragment` is the only concept you need to learn.
//
// A `*Fragment` is usually part of a SQL query, which has exactly the same bind
// var (`?` / `$n`) syntax as `database/sql`, but more than that, it allows you
// to bind other fragment builders.
//
// The `*Fragment` is usually created by `sqlf.F()`.
//
//	query, args, _ := sqlf.F(
//		"SELECT * FROM foo WHERE ?",
//		sqlf.Join(
//			" AND "
//			sqlf.F("baz = $1", true),
//			sqlf.F("bar BETWEEN ? AND ?", 1, 100),
//		),
//	).BuildQuery(syntax.Dollar)
//	fmt.Println(query)
//	fmt.Println(args)
//	// Output:
//	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
//	// [true 1 100]
package sqlf

// Builder is a builder that builds a fragment.
type Builder interface {
	// Build builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	Build(ctx *Context) (query string, err error)
}
