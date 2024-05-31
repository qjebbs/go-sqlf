package sqlb_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func TestQueryBuilder(t *testing.T) {
	var (
		users = sqlb.NewTableAliased("users", "u")
		foo   = sqlb.NewTableAliased("foo", "f")
		bar   = sqlb.NewTableAliased("bar", "b")
	)
	q := sqlb.NewQueryBuilder().
		Distinct().
		With(
			users.Name,
			sqlf.Fa("SELECT * FROM users WHERE type=$1", "user"),
		).
		With("xxx", sqlf.F("SELECT 1 AS whatever")) // should be ignored
	q.Select(foo.Columns("id", "name")...).
		From(users).
		LeftJoinOptional(foo, sqlf.Ff(
			"#f1=#f2",
			foo.Column("user_id"),
			users.Column("id"),
		)).
		LeftJoinOptional(bar, sqlf.Ff( // not referenced, should be ignored
			"#f1=#f2",
			bar.Column("user_id"),
			users.Column("id"),
		)).
		Where2(users.Column("id"), "=", 1).
		Union(
			sqlb.NewQueryBuilder().
				Select(foo.Columns("id", "name")...).
				From(foo).
				Where(sqlf.F("#f1>$1 AND #f1<$2").
					WithFragments(foo.Column("id")).
					WithArgs(10, 20),
				),
		)
	gotQuery, gotArgs, err := q.BuildQuery(syntax.Dollar)
	if err != nil {
		t.Fatal(err)
	}
	wantQuery := "With users AS (SELECT * FROM users WHERE type=$1) SELECT DISTINCT f.id, f.name FROM users AS u LEFT JOIN foo AS f ON f.user_id=u.id WHERE u.id=$2 UNION (SELECT f.id, f.name FROM foo AS f WHERE f.id>$3 AND f.id<$4)"
	wantArgs := []any{"user", 1, 10, 20}
	if wantQuery != gotQuery {
		t.Errorf("got:\n%s\nwant:\n%s", gotQuery, wantQuery)
	}
	if !reflect.DeepEqual(wantArgs, gotArgs) {
		t.Errorf("want:\n%v\ngot:\n%v", wantArgs, gotArgs)
	}
}
