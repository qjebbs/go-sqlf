package sqlb_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/sqlb"
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

func TestQueryBuilderDistinctElimination(t *testing.T) {
	var (
		users = sqlb.NewTableAliased("users", "u")
		locs  = sqlb.NewTableAliased("locs", "l")
		foo   = sqlb.NewTableAliased("foo", "f")
		bar   = sqlb.NewTableAliased("bar", "b")
	)
	q := sqlb.NewQueryBuilder().
		Distinct().
		With(sqlb.NewTableAliased("xxx", ""), sqlf.F("SELECT 1 AS whatever")). // should be ignored
		With(locs, sqlf.F("SELECT user_id AS id, loc FROM user_locs WHERE country_code = ?", "cn")).
		With(
			users,
			// CTE references another
			sqlf.F("SELECT * FROM ? AS ? INNER JOIN ? AS ? ON ?=?",
				users.Name, users.Alias, locs.Name, locs.Alias,
				users.Column("id"), locs.Column("id"),
			),
		)
	q.Select(foo.Columns("id", "name")...).
		From(users).
		LeftJoinOptional(foo, sqlf.F(
			"?=?",
			foo.Column("user_id"),
			users.Column("id"),
		)).
		LeftJoinOptional(bar, sqlf.F( // not referenced, should be ignored
			"?=?",
			bar.Column("user_id"),
			users.Column("id"),
		))
	gotQuery, gotArgs, err := q.BuildQuery(syntax.Dollar)
	if err != nil {
		t.Fatal(err)
	}
	wantQuery := "With locs AS (SELECT user_id AS id, loc FROM user_locs WHERE country_code = $1), users AS (SELECT * FROM users AS u INNER JOIN locs AS l ON u.id=l.id) SELECT DISTINCT f.id, f.name FROM users AS u LEFT JOIN foo AS f ON f.user_id=u.id"
	wantArgs := []any{"cn"}
	if wantQuery != gotQuery {
		t.Errorf("got:\n%s\nwant:\n%s", gotQuery, wantQuery)
	}
	if !reflect.DeepEqual(wantArgs, gotArgs) {
		t.Errorf("want:\n%v\ngot:\n%v", wantArgs, gotArgs)
	}
}

func TestQueryBuilderGroupbyElimination(t *testing.T) {
	var (
		foo = sqlb.NewTableAliased("foo", "f")
		bar = sqlb.NewTableAliased("bar", "b")
		baz = sqlb.NewTableAliased("baz", "z")
	)
	q := sqlb.NewQueryBuilder().
		With(
			baz,
			sqlf.F("SELECT * FROM baz WHERE type=$1", "user"),
		)
	q.Select(foo.Columns("id", "bar")...).
		From(foo).
		LeftJoinOptional(baz, sqlf.F(
			"?=?",
			foo.Column("baz_id"),
			baz.Column("id"),
		)).
		LeftJoinOptional(bar, sqlf.F( // not referenced, should be ignored
			"?=?",
			bar.Column("baz_id"),
			baz.Column("id"),
		)).
		Where2(foo.Column("id"), "=", 1).
		GroupBy(foo.Column("id"))
	gotQuery, gotArgs, err := q.BuildQuery(syntax.Dollar)
	if err != nil {
		t.Fatal(err)
	}
	wantQuery := "SELECT f.id, f.bar FROM foo AS f WHERE f.id=$1 GROUP BY f.id"
	wantArgs := []any{1}
	if wantQuery != gotQuery {
		t.Errorf("got:\n%s\nwant:\n%s", gotQuery, wantQuery)
	}
	if !reflect.DeepEqual(wantArgs, gotArgs) {
		t.Errorf("want:\n%v\ngot:\n%v", wantArgs, gotArgs)
	}
}
