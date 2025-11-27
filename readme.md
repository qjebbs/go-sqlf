## Go SQL Fragment (go-sqlf)

Package `sqlf` is dedicated to building SQL queries by composing fragments.

> `v3` makes the package exceptionally lightweight and easy to use by preserving only the native bind variable syntax of `database/sql`.

Unlike other SQL builders or ORMs, `*Fragment` is the only concept you need to understand.
It uses the same bind variable syntax (`?` / `$1`) as `database/sql`, and also supports binding other fragment builders for flexible query composition.

A `*Fragment` is usually created by `F()`.

```go
import (
	"fmt"
	"github.com/qjebbs/go-sqlf/v4"
)
func Example_basic() {
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join(
			" AND "
			sqlf.F("baz = $1", true),             
			sqlf.F("bar BETWEEN ? AND ?", 1, 100),
		),
	).BuildQuery(sqlf.BindStyleDollar)
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
