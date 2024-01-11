package sqlb_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/sqlb"
	"github.com/qjebbs/go-sqlf/syntax"
)

func TestQueryBuilder(t *testing.T) {
	var (
		users = sqlb.NewTable("users", "u")
		foo   = sqlb.NewTable("foo", "f")
		bar   = sqlb.NewTable("bar", "b")
	)
	q := sqlb.NewQueryBuilder().
		BindVar(syntax.Dollar).Distinct().
		With(users.Name, &sqlf.Fragment{
			Raw:  "SELECT * FROM users WHERE type=$1",
			Args: []any{"user"},
		}).
		With("xxx", &sqlf.Fragment{Raw: "SELECT 1 AS whatever"}) // should be ignored
	q.Select(foo.Columns("id", "name")...).
		From(users).
		LeftJoinOptional(foo, &sqlf.Fragment{
			Raw: "#c1=#c2",
			Columns: []*sqlf.TableColumn{
				foo.Column("user_id"),
				users.Column("id"),
			},
		}).
		LeftJoinOptional(bar, &sqlf.Fragment{ // not referenced, should be ignored
			Raw: "#c1=#c2",
			Columns: []*sqlf.TableColumn{
				bar.Column("user_id"),
				users.Column("id"),
			},
		}).
		Where2(users.Column("id"), "=", 1).
		Union(
			sqlb.NewQueryBuilder().
				BindVar(syntax.Dollar).
				Select(foo.Columns("id", "name")...).
				From(foo).
				Where(&sqlf.Fragment{
					Raw:     "#c1>$1 AND #c1<$2",
					Columns: foo.Columns("id"),
					Args:    []any{10, 20},
				}),
		)
	gotQuery, gotArgs, err := q.Build()
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
