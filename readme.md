## Go SQL Fragment (go-sqlf)

Package `sqlf` is dedicated to building SQL queries by composing fragments.

Unlike other SQL builders or ORMs, `*Fragment` is the only concept you need to understand.
It uses the same bind variable syntax (`?` / `$1`) as `database/sql`, which can refer to both ordinary args and fragment builders in args.

A `*Fragment` is created by `F()`.

```go
import (
	"context"
	"fmt"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/dialect"
)

func Example_basic() {
	b := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join([]sqlf.Builder{
			sqlf.F("baz = $1", true),
			sqlf.F("bar BETWEEN ? AND ?", 1, 100),
		}, " AND "),
	)
	ctx := sqlf.NewContext(context.Background(), dialect.PostgreSQL{})
	query, args, _ := b.Build(ctx)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}
```

## The Go SQL Tools Family

This project is part of a family of Go SQL tools, each designed for a different level of abstraction and automation:

1. **go-sqlf (this project)** — Minimalist SQL fragment builder. For simple, manual SQL composition with parameter binding and zero magic.
2. **[go-sqlb](https://github.com/qjebbs/go-sqlb)** — Advanced SQL builder. For programmatically building complex queries (CTE, JOIN, expressions, etc.) with chainable, declarative, and composable APIs.
3. **[go-sqlm](https://github.com/qjebbs/sqlm)** — Struct mapping. Declarative struct mapping, automatic CRUD, batch operations, and high-performance zero-reflection code generation.