## Go SQL Fragment (go-sqlf)

Package `sqlf` is dedicated to building SQL queries by composing fragments.

> `v4` makes the package exceptionally lightweight and easy to use by preserving only the native bind variable syntax of `database/sql`.

Unlike other SQL builders or ORMs, `*Fragment` is the only concept you need to understand.
It uses the same bind variable syntax (`?` / `$1`) as `database/sql`, and also supports binding other fragment builders for flexible query composition.

A `*Fragment` is created by `F()`.

```go
import (
	"fmt"
	"github.com/qjebbs/go-sqlf/v4"
)
func Example_basic() {
	ctx := sqlf.NewContext(context.Background(), dialect.PostgreSQL{})
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join(
			" AND ",
			sqlf.F("baz = $1", true),
			sqlf.F("bar BETWEEN ? AND ?", 1, 100),
		),
	).Build(ctx)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}
```

## Query Builder

Package [go-sqlb](https://github.com/qjebbs/go-sqlb) provides complex SQL builders and struct mapping capabilities, while `go-sqlf` is the underlying foundation.
