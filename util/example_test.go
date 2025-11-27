package util_test

import (
	"fmt"
	"time"

	"github.com/qjebbs/go-sqlf/v4/util"
)

func ExampleFlatten() {
	print := func(v any) {
		fmt.Printf("%#v\n", v)
	}
	print(util.Flatten(1, 2, 3))
	print(util.Flatten([]int{1, 2, 3}))
	print(util.Flatten(&[]int{1, 2, 3}))
	print(util.Flatten([3]int{1, 2, 3}))
	print(util.Flatten(1, []int{2, 3}, []string{"a", "b", "c"}))
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
	interpolated, err := util.Interpolate(query, args, util.WithInterpolateTimeFormat("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}
	fmt.Println(interpolated)
	// Output:
	// SELECT * FROM foo WHERE status = 'ok' AND created_at > '1970-01-01 08:00:00'
}
