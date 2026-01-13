package util_test

import (
	"fmt"
	"time"

	"github.com/qjebbs/go-sqlf/v4/dialect"
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
	interpolated, err := util.Interpolate(dialect.PostgreSQL{}, query, args)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(interpolated)
	// Output:
	// SELECT * FROM foo WHERE status = 'ok' AND created_at > '2026-01-13 00:00:00+00:00'
}
