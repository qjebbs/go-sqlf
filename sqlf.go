// Package sqlf focuses only on building SQL queries by combining fragments.
// Low reusability and scalability are the main challenges we face when
// writing SQL, the package is designed to solve these problems.
//
// # Fragment
//
// Unlike any other sql builder or ORMs, Fragment is the only concept
// you need to learn.
//
// Fragment is usually a part of a SQL query, for example, combining main fragment and any number of condition fragments, we can get a complete query.
//
//	query, args, _ := sqlf.Ff(
//		`SELECT * FROM foo WHERE #join('#fragment', ' AND ')`,
//		sqlf.Fa("baz = $1", true),
//		sqlf.Fa("bar BETWEEN ? AND ?)", 1, 100),
//	).BuildQuery()
//	fmt.Println(query)
//	fmt.Println(args)
//	// Output:
//	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3)
//	// [true 1 100]
//
// Explanation:
//
//   - We pay attention only to the references inside the fragment, for example, use $1 to refer Fragment.Args[0], or ? to refer Fragment.Args in order.
//   - #join, #f, etc., are preprocessing functions, which will be explained later.
//   - See Example (Basic1) for what happend inside a *sqlf.Fragment.
//
// # Preprocessing Functions
//
//   - f, fragment: fragment properties at index, e.g. #f1
//   - join: Join the template with separator, e.g. #join('#f', ', '), #join('#arg', ',', 3), #join('#arg', ',', 3, 6)
//   - arg: fragment args at index, usually used in #join().
//
// Note:
//   - #f1 is equivalent to #f(1), which is a special syntax to call preprocessing functions when an integer (usually an index) is the only argument.
//   - Expressions in the #join template are functions, not function calls.
//   - You can register custom functions to the build context, see Context.Funcs.
package sqlf

// QueryBuilder is the interface for sql builders.
type QueryBuilder interface {
	// BuildQuery builds and returns the query and args.
	BuildQuery() (query string, args []any, err error)
}

// FragmentBuilder is a builder that builds a fragment.
type FragmentBuilder interface {
	// BuildFragment builds as a fragment with the context.
	// The args should be committed to the ctx if any.
	BuildFragment(ctx *Context) (query string, err error)
}
