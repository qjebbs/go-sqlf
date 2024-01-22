package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/syntax"
)

func Example_basic() {
	query, args, _ := (&sqlf.Fragment{
		Raw: `SELECT * FROM foo WHERE #join('#fragment', ' AND ')`,
		Fragments: []*sqlf.Fragment{
			sqlf.FArgs(`bar IN (#join('#argDollar', ', '))`, 1, 2, 3),
			sqlf.FArgs(`baz = $1`, true),
		},
	}).Build()
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT * FROM foo WHERE bar IN ($1, $2, $3) AND baz = $4
	// [1 2 3 true]
}

func Example_select() {
	selects := &sqlf.Fragment{
		Raw: "SELECT #join('#column', ', ')",
	}
	from := &sqlf.Fragment{
		Raw: "FROM #t1",
	}
	where := &sqlf.Fragment{
		Prefix: "WHERE",
		Raw:    "#join('#fragment', ' AND ')",
	}
	builder := &sqlf.Fragment{
		Raw: "#join('#fragment', ' ')",
		Fragments: []*sqlf.Fragment{
			selects,
			from,
			where,
		},
	}

	var users sqlf.Table = "users"
	selects.WithColumns(users.Expressions("id", "name", "email")...)
	from.WithTables(users)
	where.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1 IN (#join('#argDollar', ', '))",
		Columns: users.Expressions("id"),
		Args:    []any{1, 2, 3},
	})
	where.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1 = $1",
		Columns: users.Expressions("active"),
		Args:    []any{true},
	})

	bulit, args, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(bulit)
	fmt.Println(args)
	// Output:
	// SELECT id, name, email FROM users WHERE id IN ($1, $2, $3) AND active = $4
	// [1 2 3 true]
}

func Example_update() {
	update := &sqlf.Fragment{
		Raw: "UPDATE #t1 SET #join('#c=#argDollar', ', ')",
	}
	where := &sqlf.Fragment{
		Prefix: "WHERE",
		Raw:    "#join('#fragment', ' AND ')",
	}
	// consider wrapping it with your own builder
	// to provide a more friendly APIs
	builder := &sqlf.Fragment{
		Raw: "#join('#fragment', ' ')",
		Fragments: []*sqlf.Fragment{
			update,
			where,
		},
	}

	var users sqlf.Table = "users"
	update.WithTables(users)
	update.WithColumns(users.Expressions("name", "email")...)
	update.WithArgs("alice", "alice@example.org")
	where.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1=$1",
		Columns: users.Expressions("id"),
		Args:    []any{1},
	})

	bulit, args, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(bulit)
	fmt.Println(args)
	// Output:
	// UPDATE users SET name=$1, email=$2 WHERE id=$3
	// [alice alice@example.org 1]
}

func Example_globalArgs() {
	// this example shows how to use Global Args by using
	// *sqlf.ArgsProperty and custom function, so that we
	// don't have to put Args into every fragment, which leads
	// to a list of redundant args.
	ids := sqlf.NewArgsProperty([]any{1, 2, 3})
	ctx := sqlf.NewContext()
	err := ctx.Funcs(sqlf.FuncMap{
		"_id": func(i int) (string, error) {
			return ids.Build(ctx, i, syntax.Dollar)
		},
	})
	if err != nil {
		panic(err)
	}
	fragment := &sqlf.Fragment{
		Raw: "#join('#fragment', ' UNION ')",
		Fragments: []*sqlf.Fragment{
			{Raw: "SELECT id, 'foo' typ, count FROM foo WHERE id IN (#join('#_id', ', '))"},
			{Raw: "SELECT id, 'bar' typ, count FROM bar WHERE id IN (#join('#_id', ', '))"},
		},
	}
	bulit, err := fragment.BuildContext(ctx)
	if err != nil {
		panic(err)
	}
	args := ctx.Args()
	fmt.Println(bulit)
	fmt.Println(args)
	// Output:
	// SELECT id, 'foo' typ, count FROM foo WHERE id IN ($1, $2, $3) UNION SELECT id, 'bar' typ, count FROM bar WHERE id IN ($1, $2, $3)
	// [1 2 3]
}
