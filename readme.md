Package `sqlf` focuses on building SQL queries by free combination of fragments. 
Low reusability and scalability are the main challenges we face when writing SQL, 
the package is designed to solve these problems.

## Fragment

Unlike any other sql builder or ORMs, `Fragment` is the only concept you need to learn.

Fragment is usually a part of a SQL query, which uses exactly the same syntax as `database/sql`, but provides the ability to combine them in any way.

```go
query, args, _ := sqlf.Ff(
	"SELECT * FROM foo WHERE #join('#fragment', ' AND ')",
	sqlf.Fa("baz = $1", true),
	sqlf.Fa("bar BETWEEN ? AND ?", 1, 100),
).BuildQuery()
fmt.Println(query)
fmt.Println(args)
// Output:
// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
// [true 1 100]
```

Explanation:

- We pay attention only to the references inside a fragment, e.g., 
use `$1` to refer `Fragment.Args[0]`, or `?` to refer `Fragment.Args` in order.
- `#join`, `#arg`, `#f`, etc., are preprocessing functions, which will be explained later.

See [example_test.go](./example_test.go) for more examples.

## Preprocessing Functions

| name        | description                      | example                |
| ----------- | -------------------------------- | ---------------------- |
| f, fragment | fragment properties at index     | #f1, #fragment1        |
| join        | Join the template with separator | #join('#f', ' AND ')   |
|             | Join from index 3 to end         | #join('#f', ',', 3)    |
|             | Join from index 3 to 6           | #join('#f', ',', 3, 6) |
| arg         | fragment args at index           | #join('#arg', ',')     |

Note:
  - #f1 is equivalent to #f(1), which is a special syntax to call preprocessing functions when an integer (usually an index) is the only argument.
  - Expressions in the #join template are functions, not function calls.

You can register custom preprocessing functions to the build context.

```go
ctx := sqlf.NewContext()
ids := sqlf.NewArgsProperties(1, 2, 3)
err := ctx.Funcs(sqlf.FuncMap{
	"_id": func(i int) (string, error) {
		return ids.Build(ctx, i)
	},
})
if err != nil {
	fmt.Println(err)
	return
}
fragment := sqlf.Ff(
	"#join('#fragment', '\nUNION\n')",
	sqlf.Fa("SELECT id, 'foo' typ, count FROM foo WHERE id IN (#join('#_id', ', '))"),
	sqlf.Fa("SELECT id, 'bar' typ, count FROM bar WHERE id IN (#join('#_id', ', '))"),
)
query, err := fragment.BuildFragment(ctx)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(query)
fmt.Println(ctx.Args())
// Output:
// SELECT id, 'foo' typ, count FROM foo WHERE id IN ($1, $2, $3)
// UNION
// SELECT id, 'bar' typ, count FROM bar WHERE id IN ($1, $2, $3)
// [1 2 3]
```

## QueryBuilder

`*sqlb.QueryBuilder` is a high-level abstraction of SQL queries for building complex queries,
with `*sqlf.Fragment` as its underlying foundation.

See [sqlb/example_test.go](./sqlb/example_test.go) for examples.
