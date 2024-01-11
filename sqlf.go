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
// inside the fragment, for example, use "$1" to refer the first element of s.Args.
//
// The syntax of the fragment is exactly the same as the syntax of the "database/sql",
// plus preprocessing functions support:
//
//	SELECT * FROM foo WHERE id IN ($1, $2, $3) AND #fragment(1)
//	SELECT * FROM foo WHERE id IN (?, ?, ?) AND #fragment(1)
//	SELECT * FROM foo WHERE #join('#fragment', ' AND ')
//
// # Preprocessing Functions
//
//   - c, col, column: Column at index, e.g. #c1, #c(1)
//   - t, table				: Table name / alias at index, e.g. #t1, #t(1)
//   - f, fragment		    : Fragment at index, e.g. #f1, #f(1)
//   - b, builder		    : Builder at index, e.g. #b1, #b(1)
//   - argDollar			: arg at index with style $x, usually used in #join().
//   - argQuestion			: arg at index with style ?, usually used in #join().
//   - ctxArgDollar 		: arg from global context with style $x, e.g.: #ctxArgDollar1, #ctxArgDollar(1)
//   - ctxArgQuestion 		: arg from global context with style ?, e.g.: #ctxArgQuestion1, #ctxArgQuestion(1)
//   - join 				: Join the template by the separator, e.g. #join('#column', ', '), #join('#argDollar', ',', 3), #join('#argDollar', ',', 3, 6)
//
// Note:
//   - References in the #join template are functions, not function calls.
//   - #c1 is equivalent to #c(1), which is a special syntax to call preprocessing functions when a number is the only argument.
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
