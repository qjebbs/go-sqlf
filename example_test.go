package sqlf_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf"
)

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
		Prefix: "",
		Raw:    "UPDATE #t1 SET #join('#c=#argDollar', ', ')",
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

func ExampleContext_Funcs() {
	ctx := sqlf.NewContext()
	ctx.Funcs(sqlf.FuncMap{
		"interpolate": func(ctx *sqlf.FragmentContext, i int) (string, error) {
			// avoid "arg i is not used" error
			ctx.ReportUsedArg(i)
			return fmt.Sprint(ctx.This.Args[i-1]), nil
		},
	})
	fragment := &sqlf.Fragment{
		Raw:  "#interpolate(1)",
		Args: []any{1},
	}
	bulit, err := fragment.BuildContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(bulit)
	fmt.Println(ctx.BuiltArgs())
	// Output:
	// 1
	// []
}
