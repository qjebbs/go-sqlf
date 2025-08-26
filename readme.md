Package `sqlf` focuses on building SQL queries by free combination of fragments. 

The package exports only a few functions and methods, but improves a lot on the 
reusability and extensibility of SQL, which are the main challenges we encounter 
when writing SQL.

## Fragment

Unlike any other sql builder or ORMs, `*Fragment` is the only concept you need to learn.

A `*Fragment` is usually part of a SQL query, which has exactly the same bind 
var syntax (`?` / `$1`) as `database/sql`, but more than that, it allows you 
to bind other fragment builders.

The `*Fragment` is usually created by `F()`.

```go
import (
	"fmt"
	"github.com/qjebbs/go-sqlf/v3"
)
func Example_basic() {
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join(
			" AND "
			sqlf.F("baz = $1", true),             
			sqlf.F("bar BETWEEN ? AND ?", 1, 100),
		),
	).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}
```

## QueryBuilder

Package sqlb provides a complex SQL query builder shipped  with WITH-CTE / JOIN 
Elimination capabilities, while `*sqlf.Fragment` is the underlying foundation.

See [sqlb/example_test.go](./sqlb/example_test.go) for examples.
