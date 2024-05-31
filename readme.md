Package `sqlf` focuses on building SQL queries by free combination of fragments. 

The package exports only a few functions and methods, but improves a lot on the 
reusability and extensibility of SQL, which are the main challenges we encounter 
when writing SQL.

## Fragment

Unlike any other sql builder or ORMs, `Fragment` is the only concept you need to learn.

Fragment is usually a part of a SQL query, which uses exactly the same syntax as 
`database/sql`, but provides the ability to combine them in any way.

```go
import (
	"fmt"
	"github.com/qjebbs/go-sqlf/v2"
)
func Example_basic() {
	query, args, _ := sqlf.Ff(
		"SELECT * FROM foo WHERE #join('#fragment', ' AND ')", // join fragments
		sqlf.Fa("baz = $1", true),                             // `database/sql` style
		sqlf.Fa("bar BETWEEN ? AND ?", 1, 100),                // `database/sql` style
	).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}
```

Explanation:

- We pay attention only to the references inside a fragment, not between fragments.
- `#join`, `#arg`, `#f`, etc., are preprocessing functions, which will be explained later.
- See `Example_deeperLook` of [example_test.go](./example_test.go) for what happend inside the *sqlf.Fragment.

## Preprocessing Functions

| name        | description                      | example                |
| ----------- | -------------------------------- | ---------------------- |
| f, fragment | fragments at index               | #f1, #fragment1        |
| join        | Join the template with separator | #join('#f', ' AND ')   |
|             | Join from index 3 to end         | #join('#f', ',', 3)    |
|             | Join from index 3 to 6           | #join('#f', ',', 3, 6) |
| arg         | arguments at index               | #join('#arg', ',')     |

Note:
  - #f1 is equivalent to #f(1), which is a special syntax to call preprocessing functions when an integer (usually an index) is the only argument.
  - Expressions in the #join template are functions, not function calls.

See Example `ContextWithFuncs` of [example_test.go](./example_test.go) for how to 
register custom preprocessing functions, and implementing global arguments/fragments.

## QueryBuilder

`*sqlb.QueryBuilder` is a high-level abstraction of SQL queries for building complex queries,
with `*sqlf.Fragment` as its underlying foundation.

See [sqlb/example_test.go](./sqlb/example_test.go) for examples.
