package sqlb_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/sqlb"
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

func ExampleQueryBuilder_BuildQuery() {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.NewTableAliased("bar", "b")
	)
	b := sqlb.NewQueryBuilder().
		Select(foo.Column("*")).
		From(foo).
		InnerJoin(bar, sqlf.F(
			"?=?",
			bar.Column("foo_id"),
			foo.Column("id"),
		)).
		Where(sqlf.F(
			"($2=$1 OR $3=$1)",
			1, foo.Column("a"), foo.Column("b"),
		)).
		Where2(bar.Column("c"), "=", 2)

	query, args, err := b.BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	query, args, err = b.BuildQuery(syntax.Question)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.* FROM foo AS f INNER JOIN bar AS b ON b.foo_id=f.id WHERE (f.a=$1 OR f.b=$1) AND b.c=$2
	// [1 2]
	// SELECT f.* FROM foo AS f INNER JOIN bar AS b ON b.foo_id=f.id WHERE (f.a=? OR f.b=?) AND b.c=?
	// [1 1 2]
}

func ExampleQueryBuilder_LeftJoinOptional() {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.NewTableAliased("bar", "b")
	)
	query, args, err := sqlb.NewQueryBuilder().
		Distinct(). // *QueryBuilder trims optional joins only when SELECT DISTINCT is used.
		Select(foo.Column("*")).
		From(foo).
		// declare an optional LEFT JOIN
		LeftJoinOptional(bar, sqlf.F(
			"?=?",
			bar.Column("foo_id"),
			foo.Column("id"),
		)).
		// don't touch any columns of "bar", so that it can be trimmed
		Where2(foo.Column("id"), ">", 1).
		BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT DISTINCT f.* FROM foo AS f WHERE f.id>$1
	// [1]
}

func ExampleQueryBuilder_With() {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.NewTableAliased("bar", "b")
		cte = sqlb.NewTableAliased("bar_type_1", "b1")
	)
	query, args, err := sqlb.NewQueryBuilder().
		With(
			cte,
			sqlf.F(
				"SELECT * FROM ? AS ? WHERE ?=?",
				bar.Name, bar.Alias, bar.Column("type"), 1,
			)).
		Select(
			foo.Column("*"),
			cte.Column("*"),
		).
		From(foo).
		LeftJoinOptional(cte, sqlf.F(
			"?=?",
			cte.Column("foo_id"),
			foo.Column("id"),
		)).
		BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// With bar_type_1 AS (SELECT * FROM bar AS b WHERE b.type=$1) SELECT f.*, b1.* FROM foo AS f LEFT JOIN bar_type_1 AS b1 ON b1.foo_id=f.id
	// [1]
}

func ExampleQueryBuilder_Union() {
	var foo = sqlb.NewTableAliased("foo", "f")
	column := foo.Column("*")
	query, args, err := sqlb.NewQueryBuilder().
		Select(column).
		From(foo).
		Where2(foo.Column("id"), " = ", 1).
		Union(
			sqlb.NewQueryBuilder().
				From(foo).
				WhereIn(foo.Column("id"), []any{2, 3, 4}).
				Select(column),
		).
		BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.* FROM foo AS f WHERE f.id = $1 UNION (SELECT f.* FROM foo AS f WHERE f.id IN ($2, $3, $4))
	// [1 2 3 4]
}

func ExampleNoDeps() {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.Table("bar")
	)
	q := sqlb.NewQueryBuilder().
		Select(foo.Column("bar")).
		From(foo).
		Where(
			// will not report 'b' (table 'bar') undefined
			sqlf.F(
				"? IN (?)",
				foo.Column("id"),
				sqlb.NoDeps(sqlf.F(
					"SELECT id FROM ?", bar,
				)),
			),
		)
	query, args, err := q.BuildQuery(syntax.Dollar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.bar FROM foo AS f WHERE f.id IN (SELECT id FROM bar)
	// []
}
