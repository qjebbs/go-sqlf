package sqlb_test

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func ExampleQueryBuilder_BuildQuery() {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.NewTableAliased("bar", "b")
	)
	b := sqlb.NewQueryBuilder().
		Select(foo.Column("*")).
		From(foo).
		InnerJoin(bar, sqlf.Ff(
			"#f1=#f2",
			bar.Column("foo_id"),
			foo.Column("id"),
		)).
		Where(sqlf.F("(#f1=$1 OR #f2=$1)").
			WithFragments(foo.Column("a"), foo.Column("b")).
			WithArgs(1),
		).
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
		LeftJoinOptional(bar, sqlf.Ff(
			"#f1=#f2",
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
			cte.Name,
			sqlf.F("SELECT * FROM #f1 AS #f2 WHERE #f3=$1").
				WithFragments(
					bar.Name, bar.Alias,
					bar.Column("type"),
				).
				WithArgs(1)).
		Select(
			foo.Column("*"),
			cte.Column("*"),
		).
		From(foo).
		LeftJoinOptional(cte, sqlf.Ff(
			"#f1=#f2",
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
