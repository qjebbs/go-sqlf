package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/util"
)

func Example_basic() {
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join(
			" AND ",
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

func Example_insert() {
	var table = sqlf.F("users")
	var fields = []any{
		sqlf.F("name"),
		sqlf.F("email"),
	}
	var values = [][]any{
		{"alice", "alice@example.org"},
		{"bob", "bob@example.org"},
	}

	f := sqlf.F(
		"INSERT INTO ? (?) VALUES ?",
		table,
		sqlf.Join(", ", fields...),
		sqlf.Join(", ", util.Map(values, func(value []any) any {
			return sqlf.F(
				"(?)",
				sqlf.Join(", ", value...),
			)
		})...),
	)

	query, args, err := f.BuildQuery(sqlf.BindStyleDollar)
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
