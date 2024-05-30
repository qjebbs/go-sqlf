package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func Example_basic1() {
	// This example is equivalent to Exmaple Basic2 (which is more concise), but
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
		Raw:       "SELECT * FROM foo WHERE #join('#fragment', ' AND ')",
		Fragments: []sqlf.FragmentBuilder{a, b},
	}).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}
func Example_basic2() {
	query, args, _ := sqlf.Ff(
		"SELECT * FROM foo WHERE #join('#fragment', ' AND ')",
		sqlf.Fa("baz = $1", true),
		sqlf.Fa("bar BETWEEN ? AND ?", 1, 100),
	).BuildQuery(syntax.Dollar)
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE baz = $1 AND bar BETWEEN $2 AND $3
	// [true 1 100]
}

func Example_select() {
	selects := sqlf.Ff("SELECT #join('#fragment', ', ')")
	from := sqlf.Ff("FROM #f1")
	where := sqlf.Ff("#join('#fragment', ' AND ')").WithPrefix("WHERE")
	builder := sqlf.Ff("#join('#fragment', ' ')", selects, from, where)

	var users sqlb.Table = "users"
	selects.WithFragments(
		users.AnonymousColumn("id"),
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	)
	from.WithFragments(users)
	where.WithFragments(
		sqlf.F("#f1 IN (#join('#arg', ', '))").
			WithFragments(users.AnonymousColumn("id")).
			WithArgs(1, 2, 3),
	)
	where.AppendFragments(
		sqlf.F("#f1 = $1").
			WithFragments(users.AnonymousColumn("active")).
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
	update := sqlf.Fa("UPDATE #f1")
	fieldValues := sqlf.Fa("SET #join('#fragment=#arg', ', ')")
	where := sqlf.Fa("#join('#fragment', ' AND ')").WithPrefix("WHERE")
	builder := sqlf.Ff("#join('#fragment', ' ')", update, fieldValues, where)

	var users sqlb.Table = "users"
	update.WithFragments(users)
	fieldValues.WithFragments(
		users.AnonymousColumn("name"),
		users.AnonymousColumn("email"),
	)
	fieldValues.WithArgs("alice", "alice@example.org")
	where.AppendFragments(
		sqlf.F("#f1=$1").
			WithFragments(users.AnonymousColumn("id")).
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

func Example_unalignedJoin() {
	// this example demonstrates how #join() works between unaligned properties.
	// it leaves the extra property items (.Args[2:] here) unused, which leads to an error.
	// to make it work, we use #noUnusedError() to suppress the error.
	ctx, err := sqlf.ContextWithFuncs(sqlf.NewContext(syntax.Dollar), sqlf.FuncMap{
		"noUnusedError": func(ctx *sqlf.Context) {
			args := ctx.Fragment().Args
			for i := 2; i < len(args); i++ {
				args[i].ReportUsed()
			}
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	foo := sqlb.Table("foo")
	b := sqlf.F("#noUnusedError() UPDATE foo SET #join('#fragment=#arg', ', ')").
		WithFragments(foo.AnonymousColumn("bar"), foo.AnonymousColumn("baz")).
		WithArgs(1, 2, 3, true, false)
	query, err := b.BuildFragment(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(ctx.Args())
	// Output:
	// UPDATE foo SET bar=$1, baz=$2
	// [1 2]
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
}
