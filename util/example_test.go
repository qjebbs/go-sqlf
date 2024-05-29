package util_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/util"
)

func ExampleArgs() {
	print := func(v any) {
		fmt.Printf("%#v\n", v)
	}
	print(util.ArgsFlatted(1, 2, 3))
	print(util.ArgsFlatted([]int{1, 2, 3}))
	print(util.ArgsFlatted(&[]int{1, 2, 3}))
	print(util.ArgsFlatted([3]int{1, 2, 3}))
	print(util.ArgsFlatted(1, []int{2, 3}, []string{"a", "b", "c"}))
	// Output:
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3, "a", "b", "c"}
}

func ExampleInterpolate() {
	query := "SELECT * FROM foo WHERE status = ? AND created_at > ?"
	args := []any{"ok", time.Unix(0, 0)}
	interpolated, err := util.Interpolate(query, args, util.WithTimeFormat("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}
	fmt.Println(interpolated)
	// Output:
	// SELECT * FROM foo WHERE status = 'ok' AND created_at > '1970-01-01 08:00:00'
}

func ExampleCountBuilder() {
	var (
		db      *sql.DB
		builder sqlf.QueryBuilder
	)
	if db != nil && builder != nil {
		count, err := util.CountBuilder(db, builder)
		if err != nil {
			panic(err)
		}
		fmt.Println(count)
	}
	// Output:
	//
}

func ExampleCount() {
	var db *sql.DB
	if db != nil {
		count, err := util.Count(db, "SELECT * FROM foo", nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(count)
	}
	// Output:
	//
}

func ExampleScanBuilder() {
	type foo struct {
		ID   int64
		Name string
	}
	var db *sql.DB
	if db != nil {
		builder := sqlf.Fa(
			"SELECT id, name FROM foo WHERE id IN (#join('#arg', ', '))",
			1, 2, 3,
		)
		r, err := util.ScanBuilder(
			db, builder,
			func() (*foo, []any) {
				r := &foo{}
				return r, []any{&r.ID, &r.Name}
			},
		)
		if err != nil {
			panic(err)
		}
		fmt.Println(r)
	}
	// Output:
	//
}

func ExampleScan() {
	type foo struct {
		ID   int64
		Name string
	}
	var db *sql.DB
	if db != nil {
		query := "SELECT id, name FROM foo LIMIT 10"
		r, err := util.Scan(
			db, query, nil,
			func() (*foo, []any) {
				r := &foo{}
				return r, []any{&r.ID, &r.Name}
			},
		)
		if err != nil {
			panic(err)
		}
		fmt.Println(r)
	}
	// Output:
	//
}
