package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func Example_basic() {
	query, args, _ := sqlf.F(
		"SELECT * FROM foo WHERE #join('#fragment', ' AND ')", // join fragments
		sqlf.F("baz = $1", true),                              // `database/sql` style
		sqlf.F("bar BETWEEN ? AND ?", 1, 100),                 // `database/sql` style
	).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}

func Example_deeperLook() {
	// This example is equivalent to Exmaple Basic (which is more concise), but
	// it reveales what happend inside a *sqlf.Fragment.

	// *sqlf.Fragment has two types of properties storage, .Args and .Fragments.
	// Raw query can reference the contents of .Args, just like `database/sql`.
	a := &sqlf.Fragment{
		Raw:  "baz = $1",
		Args: []any{true},
	}
	b := &sqlf.Fragment{
		Raw:  "bar BETWEEN ? AND ?",
		Args: []any{1, 100},
	}
	query, args, _ := (&sqlf.Fragment{
		// Similarly, referencing .Fragments results fragments combinations.
		Raw:  "SELECT * FROM foo WHERE #join('#fragment', ' AND ')",
		Args: []any{a, b},
	}).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}

func Example_select() {
	selects := sqlf.F("SELECT #join('#fragment', ', ')")
	from := sqlf.F("FROM #f1")
	where := sqlf.F("#join('#fragment', ' AND ')").WithPrefix("WHERE")
	builder := sqlf.F("#join('#fragment', ' ')", selects, from, where)

	var users sqlb.Table = "users"
	selects.WithArgs(
		users.AnonymousColumn("id"),
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	)
	from.WithArgs(users)
	where.WithArgs(
		sqlf.F("#f1 IN (#join('#arg', ', '))").
			WithArgs(users.AnonymousColumn("id")).
			WithArgs(1, 2, 3),
	)
	where.AppendArgs(
		sqlf.F("#f1 = $1").
			WithArgs(users.AnonymousColumn("active")).
			WithArgs(true),
	)

	query, args, err := builder.BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT id, name, email FROM users WHERE id IN ($1, $2, $3) AND active = $4
	// [1 2 3 true]
}

func Example_update() {
	// consider wrapping it with your own builder to provide a more friendly APIs
	update := sqlf.F("UPDATE #f1")
	fieldValues := sqlf.F("SET #join('#fragment=#arg', ', ')")
	where := sqlf.F("#join('#fragment', ' AND ')").WithPrefix("WHERE")
	builder := sqlf.F("#join('#fragment', ' ')", update, fieldValues, where)

	var users sqlb.Table = "users"
	update.WithArgs(users)
	fieldValues.WithArgs(
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	)
	fieldValues.WithArgs("alice", "alice@example.org")
	where.AppendArgs(
		sqlf.F("#f1=$1").
			WithArgs(users.AnonymousColumn("id")).
			WithArgs(1),
	)

	query, args, err := builder.BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// UPDATE users SET name=$1, email=$2 WHERE id=$3
	// [alice alice@example.org 1]
}

func ExampleContextWithFuncs() {
	// this example shows how to use Global Args by using
	// sqlf.NewArgsProperties and custom function, so that we
	// don't have to put Args into every fragment, which leads
	// to a list of redundant args.
	ids := sqlf.NewArgsProperties(1, 2, 3)
	ctx, err := sqlf.ContextWithFuncs(sqlf.NewContext(syntax.Dollar), sqlf.FuncMap{
		"_id": func(ctx *sqlf.Context, i int) (string, error) {
			return ids.Build(ctx, i)
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fragment := sqlf.F(
		"#join('#fragment', '\nUNION\n')",
		sqlf.F("SELECT id, 'foo' typ, count FROM foo WHERE id IN (#join('#_id', ', '))"),
		sqlf.F("SELECT id, 'bar' typ, count FROM bar WHERE id IN (#join('#_id', ', '))"),
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
}
