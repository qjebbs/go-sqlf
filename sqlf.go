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
//	query, args, _ := (&sqlf.Fragment{
//	    Raw: `SELECT * FROM foo WHERE #join('#fragment', ' AND ')`,
//	    Fragments: []*sqlf.Fragment{
//	        sqlf.Fa(`bar IN (#join('#argDollar', ', '))`, 1, 2, 3),
//	        sqlf.Fa(`baz = $1`, true),
//	    },
//	}).Build()
//	fmt.Println(query)
//	fmt.Println(args)
//	// Output:
//	// SELECT * FROM foo WHERE bar IN ($1, $2, $3) AND baz = $4
//	// [1 2 3 true]
//
// Explanation:
//
//   - We pay attention only to the references inside the fragment, for example, use $1 to refer Fragment.Args[0], or ? to refer Fragment.Args in order.
//   - #join, #column, #fragment, etc., are preprocessing functions, which will be explained later.
//
// # Preprocessing Functions
//
//   - c, column: Fragment.Columns at index, e.g. #c1
//   - t, table: Fragment.Tables at index, e.g. #t1
//   - fragment: Fragment.Fragments at index, e.g. #fragment1
//   - builder: Fragment.Builders at index, e.g. #builder1
//   - argDollar: Fragment.Args at index with style $, usually used in #join().
//   - argQuestion: Fragment.Args at index with style ?, usually used in #join().
//   - join: Join the template with separator, e.g. #join('#column', ', '), #join('#argDollar', ',', 3), #join('#argDollar', ',', 3, 6)
//
// Note:
//   - #c1 is equivalent to #c(1), which is a special syntax to call preprocessing functions when a number is the only argument.
//   - Expressions in the #join template are functions, not function calls.
//   - You can register custom functions to the build context, see Context.Funcs.
package sqlf

// Builder is the interface for sql builders.
type Builder interface {
	// Build builds and returns the query and args.
	Build() (query string, args []any, err error)
	// BuildContext builds the query with the context.
	// The built args should be committed to the context, which can be
	// retrieved after building.
	BuildContext(ctx *Context) (query string, err error)
}
