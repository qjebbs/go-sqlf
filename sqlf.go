// Package sqlf focuses only on building SQL queries by combining fragments.
//
// The package exports only a few functions and methods, but improves a lot on the
// reusability and extensibility of SQL, which are the main challenges we encounter
// when writing SQL.
//
// # Fragment
//
// Unlike any other sql builder or ORMs, Fragment is the only concept
// you need to learn.
//
// Fragment is usually a part of a SQL query, which uses exactly the same syntax
// as `database/sql`, but provides the ability to combine them in any way.
//
//	query, args, _ := sqlf.Ff(
//		"SELECT * FROM foo WHERE #join('#fragment', ' AND ')", // join fragments
//		sqlf.Fa("baz = $1", true),                             // `database/sql` style
//		sqlf.Fa("bar BETWEEN ? AND ?", 1, 100),                // `database/sql` style
//	).BuildQuery(syntax.Dollar)
//	fmt.Println(query)
//	fmt.Println(args)
//	// Output:
//	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
//	// [true 1 100]
//
// Explanation:
//
//   - We pay attention only to the references inside a fragment, not between fragments.
//   - #join, #f, etc., are preprocessing functions, which will be explained later.
//   - See Example (DeeperLook) for what happend inside the *sqlf.Fragment.
//
// # Preprocessing Functions
//
//   - f, fragment: fragments at index, e.g. #f1
//   - join: Join the template with separator, e.g. #join('#f', ', '), #join('#arg', ',', 3), #join('#arg', ',', 3, 6)
//   - arg: arguments at index, usually used in #join().
//
// Note:
//   - #f1 is equivalent to #f(1), which is a special syntax to call preprocessing functions when an integer (usually an index) is the only argument.
//   - Expressions in the #join template are functions, not function calls.
//   - You can register custom functions to the build context, see ContextWithFuncs.
package sqlf

import "github.com/qjebbs/go-sqlf/v2/syntax"

// QueryBuilder is the interface for sql builders.
type QueryBuilder interface {
	// BuildQuery builds and returns the query and args.
	BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error)
}

// FragmentBuilder is a builder that builds a fragment.
type FragmentBuilder interface {
	// BuildFragment builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	BuildFragment(ctx *Context) (query string, err error)
}
