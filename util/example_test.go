package util_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/util"
)

func ExampleArgs() {
	print := func(v any) {
		fmt.Printf("%#v\n", v)
	}
	print(util.Args([]any{1, 2, 3}))
	print(util.Args([]int{1, 2, 3}))
	print(util.Args(&[]int{1, 2, 3}))
	print(util.Args([3]int{1, 2, 3}))
	print(util.Args(&[]time.Duration{time.Millisecond}))
	print(util.Args(1))
	// Output:
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1000000}
	// []interface {}{1}
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
		builder sqlf.Builder
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
		builder := &sqlf.Fragment{
			Raw:  "SELECT id, name FROM foo WHERE id IN (#join('#argDollar', ', '))",
			Args: []any{1, 2, 3},
		}
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
