package sqlf_test

import (
	"context"
	"fmt"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/util"
)

func Example_basic() {
	ctx := sqlf.ContextWithDialect(
		context.Background(),
		dialect.PostgreSQL{},
	)
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

func Example_insert() {
	var table = sqlf.F("users")
	var fields = []sqlf.Builder{
		sqlf.F("name"),
		sqlf.F("email"),
	}
	var values = [][]string{
		{"alice", "alice@example.org"},
		{"bob", "bob@example.org"},
	}

	f := sqlf.F(
		"INSERT INTO ? (?) VALUES ?",
		table,
		sqlf.Join(", ", fields...),
		sqlf.Join(", ", util.Map(values, func(value []string) sqlf.Builder {
			return sqlf.F(
				"(?)",
				sqlf.JoinArgs(", ", value...),
			)
		})...),
	)

	ctx := sqlf.ContextWithDialect(
		context.Background(),
		dialect.PostgreSQL{},
	)
	query, args, err := f.Build(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// INSERT INTO users (name, email) VALUES ($1, $2), ($3, $4)
	// [alice alice@example.org bob bob@example.org]
}
