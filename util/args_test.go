package util_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/util"
)

func TestArgs(t *testing.T) {
	type str string
	strA := str("a")
	testCases := []struct {
		args []any
		want []any
	}{
		{args: []any{1}, want: []any{1}},
		{args: []any{[]int{1, 2, 3}}, want: []any{1, 2, 3}},
		{args: []any{[]string{"a", "b", "c"}}, want: []any{"a", "b", "c"}},
		{args: []any{[]str{"a", "b", "c"}}, want: []any{str("a"), str("b"), str("c")}},
		{args: []any{[]*str{&strA}}, want: []any{&strA}},
		{args: []any{[]any{1, "a", 2, "b", 3, "c"}}, want: []any{1, "a", 2, "b", 3, "c"}},
		{
			args: []any{
				0,
				[]int{1, 2, 3},
				[]string{"a", "b", "c"},
			},
			want: []any{
				0,
				1, 2, 3,
				"a", "b", "c",
			},
		},
	}
	for _, tc := range testCases {
		got := util.Args(tc.args...)
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("want: %s, got: %s", tc.want, got)
		}
	}
}
