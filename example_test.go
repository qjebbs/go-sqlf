package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/sqlb"
	"github.com/qjebbs/go-sqlf/v3/syntax"
	"github.com/qjebbs/go-sqlf/v3/util"
)

func Example_basic() {
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE ?",
		sqlf.Join(
			" AND ",
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

func Example_select() {
	var users sqlb.Table = "users"
	selects := []sqlf.Builder{
		users.AnonymousColumn("id"),
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	}
	from := users
	where := []sqlf.Builder{
		sqlf.F(
			"? = ?",
			users.AnonymousColumn("active"), true,
		),
		sqlf.F(
			"? IN (?)",
			users.AnonymousColumn("id"),
			sqlf.Join(", ", 1, 2, 3),
		),
	}

	builder := sqlf.F(
		"SELECT ? FROM ? WHERE ?",
		sqlf.Join(", ", util.Ttoa(selects)...),
		from,
		sqlf.Join(" AND ", util.Ttoa(where)...),
	)

	query, args, err := builder.BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT id, name, email FROM users WHERE active = $1 AND id IN ($2, $3, $4)
	// [true 1 2 3]
}

func Example_insert() {
	var users sqlb.Table = "users"
	var fields = []sqlf.Builder{
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	}
	var values = [][]any{
		{"alice", "alice@example.org"},
		{"bob", "bob@example.org"},
	}

	builder := &insertBuilder{
		table:  users,
		fields: fields,
		values: values,
	}

	query, args, err := builder.BuildQuery(syntax.Dollar)
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

var _ sqlb.Builder = (*insertBuilder)(nil)

type insertBuilder struct {
	table  sqlb.Table
	fields []sqlf.Builder
	values [][]any
}

func (b *insertBuilder) BuildQuery(bindVarStyle syntax.BindVarStyle) (string, []any, error) {
	valueFragments := make([]any, len(b.values))
	for i, value := range b.values {
		valueFragments[i] = sqlf.F(
			"(?)",
			sqlf.Join(", ", value...),
		)
	}
	f := sqlf.F(
		"INSERT INTO ? (?) VALUES ?",
		b.table,
		sqlf.Join(", ", util.Ttoa(b.fields)...),
		sqlf.Join(", ", valueFragments...),
	)
	return f.BuildQuery(bindVarStyle)
}
