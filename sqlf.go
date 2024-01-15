// Package sqlf focuses only on bulding SQL queries by free combination
// of fragments. Thus, it works naturally with all sql dialects without
// having to deal with the differences between them. Unlike any other
// sql builder or ORMs, Fragment is the only concept you need to learn.
//
// # Fragment
//
// Fragment is the builder for a part of or even the full query, it allows you
// to write and combine fragments with freedom.
//
// With the help of Fragment, we pay attention only to the reference relationships
// inside the fragment, for example, use $1 to refer Fragment.Args[0].
//
// The syntax of the fragment is exactly the same as the syntax of the "database/sql",
// plus preprocessing functions support:
//
//	SELECT * FROM foo WHERE id IN ($1, $2, $3) AND #fragment(1)
//	SELECT * FROM foo WHERE id IN (?, ?, ?) AND #fragment(1)
//	SELECT * FROM foo WHERE #join('#fragment', ' AND ')
//
// Explanation:
//   - $1, $2, $3 means the first, second, third argument of the Fragment.Args.
//   - ? means the argument of the Fragment.Args in order.
//   - #fragment is a preprocessing function, which will be explained later.
//
// # Preprocessing Functions
//
//   - c, column		: Fragment.Columns at index, e.g. #c1
//   - t, table			: Fragment.Tables at index, e.g. #t1
//   - fragment		   	: Fragment.Fragments at index, e.g. #fragment1
//   - builder		   	: Fragment.Builders at index, e.g. #builder1
//   - argDollar		: Fragment.Args at index with style $, usually used in #join().
//   - argQuestion		: Fragment.Args at index with style ?, usually used in #join().
//   - globalArgDollar 	: Arg from global context with style $, e.g.: #globalArgDollar1
//   - globalArgQuestion 	: Arg from global context with style ?, e.g.: #globalArgQuestion1
//   - join 			: Join the template with separator, e.g. #join('#column', ', '), #join('#argDollar', ',', 3), #join('#argDollar', ',', 3, 6)
//
// Note:
//   - #c1 is equivalent to #c(1), which is a special syntax to call preprocessing functions when a number is the only argument.
//   - Expressions in the #join template are functions, not function calls.
package sqlf

// Builder is the interface for sql builders.
//
// Implementations tip: if a builder want to render an Arg in the *Context,
// use ctx.Arg(index).
type Builder interface {
	// Build builds and returns the query and args.
	Build() (query string, args []any, err error)
	// BuildContext builds the query with the context.
	BuildContext(ctx *Context) (query string, err error)
}
