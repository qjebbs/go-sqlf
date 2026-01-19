package util_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/qjebbs/go-sqlf/v4/util"
)

func ExampleFlattenArgs() {
	print := func(v any) {
		fmt.Printf("%#v\n", v)
	}
	print(util.FlattenArgs(1, 2, 3))
	print(util.FlattenArgs([]int{1, 2, 3}))
	print(util.FlattenArgs(&[]int{1, 2, 3}))
	print(util.FlattenArgs([3]int{1, 2, 3}))
	print(util.FlattenArgs[any](1, []int{2, 3}, []string{"a", "b", "c"}))
	// Output:
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3}
	// []interface {}{1, 2, 3, "a", "b", "c"}
}

func ExampleInterpolate() {
	query := "SELECT * FROM foo WHERE status = ? AND created_at > ?"
	args := []any{"ok", time.Date(2026, 01, 13, 0, 0, 0, 0, time.UTC)}
	interpolated, ok := util.Interpolate(query, args)
	fmt.Println(interpolated)
	fmt.Println(ok)
	// Output:
	// SELECT * FROM foo WHERE status = 'ok' AND created_at > '2026-01-13 00:00:00'
	// true
}

func ExampleInterpolate_named() {
	query := "SELECT * FROM foo WHERE status = @status AND created_at > @created_at"
	args := []any{
		sql.Named("status", "ok"),
		sql.Named(
			"created_at",
			time.Date(2026, 01, 13, 0, 0, 0, 0, time.UTC),
		),
	}
	interpolated, ok := util.Interpolate(query, args)
	fmt.Println(interpolated)
	fmt.Println(ok)
	// Output:
	// SELECT * FROM foo WHERE status = 'ok' AND created_at > '2026-01-13 00:00:00'
	// true
}
