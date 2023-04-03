package sqlb_test

import (
	"fmt"

	"git.qjebbs.com/jebbs/go-sqls"
	"git.qjebbs.com/jebbs/go-sqls/sqlb"
)

func ExampleQueryBuilder_Build() {
	foo := sqls.Table{"foo", "f"}
	bar := sqls.Table{"bar", "b"}
	query, args, err := sqlb.NewQueryBuilder(nil).
		Select(foo.AnyColumns()).
		From(foo).
		InnerJoin(bar, &sqls.Segment{
			Raw: "#c1=#c2",
			Columns: []*sqls.TableColumn{
				bar.Column("foo_id"),
				foo.Column("id"),
			},
		}).
		Where(&sqls.Segment{
			Raw:     "(#c1=$1 OR #c2=$2)",
			Columns: foo.Columns("a", "b"),
			Args:    []any{1, 2},
		}).
		Where2(bar.Column("c"), "=", 3).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.* FROM foo AS f INNER JOIN bar AS b ON b.foo_id=f.id WHERE (f.a=$1 OR f.b=$2) AND b.c=$3
	// [1 2 3]
}

func ExampleQueryBuilder_LeftJoinOptional() {
	foo := sqls.Table{"foo", "f"}
	bar := sqls.Table{"bar", "b"}

	query, args, err := sqlb.NewQueryBuilder(nil).
		Distinct(). // *QueryBuilder trims optional joins only when SELECT DISTINCT is used.
		Select(foo.AnyColumns()).
		From(foo).
		// declare an optional LEFT JOIN
		LeftJoinOptional(bar, &sqls.Segment{
			Raw: "#c1=#c2",
			Columns: []*sqls.TableColumn{
				bar.Column("foo_id"),
				foo.Column("id"),
			},
		}).
		// don't touch any columns of "bar", so that it can be trimmed
		Where2(foo.Column("id"), ">", 1).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT DISTINCT f.* FROM foo AS f WHERE f.id>$1
	// [1]
}

func ExampleQueryBuilder_With() {
	foo := sqls.Table{"foo", "f"}
	bar := sqls.Table{"bar", "b"}
	cte := sqls.Table{"bar_type_1", "b1"}
	query, args, err := sqlb.NewQueryBuilder(nil).
		With(cte, &sqls.Segment{
			Raw:     "SELECT * FROM #tAs1 WHERE #c1=$1",
			Columns: bar.Columns("type"),
			Args:    []any{1},
		}).
		Select([]*sqls.TableColumn{
			foo.AnyColumn(),
			cte.AnyColumn(),
		}).
		From(foo).
		LeftJoin(cte, &sqls.Segment{
			Raw: "#c1=#c2",
			Columns: []*sqls.TableColumn{
				cte.Column("foo_id"),
				foo.Column("id"),
			},
		}).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// With bar_type_1 AS (SELECT * FROM bar AS b WHERE b.type=$1) SELECT f.*, b1.* FROM foo AS f LEFT JOIN bar_type_1 AS b1 ON b1.foo_id=f.id
	// [1]
}

func ExampleQueryBuilder_Union() {
	foo := sqls.Table{"foo", "f"}
	columns := foo.AnyColumns()
	query, args, err := sqlb.NewQueryBuilder(nil).
		Select(columns).
		From(foo).
		Where2(foo.Column("id"), " = ", 1).
		Union(
			sqlb.NewQueryBuilder(nil).
				From(foo).
				WhereIn(foo.Column("id"), []any{2, 3, 4}).
				Select(columns),
		).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.* FROM foo AS f WHERE f.id = $1 UNION (SELECT f.* FROM foo AS f WHERE f.id IN ($2, $3, $4))
	// [1 2 3 4]
}
