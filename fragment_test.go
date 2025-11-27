package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v4"
)

func TestBuildFragment(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		style    sqlf.BindStyle
		fragment sqlf.Builder
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "build nil fragment",
			fragment: (*sqlf.Fragment)(nil),
			want:     "",
			wantArgs: nil,
		},
		{
			name:     "fragment arg",
			fragment: sqlf.F("WHERE 1=1 ?", sqlf.F("")),
			want:     "WHERE 1=1",
			wantArgs: nil,
		},
		{
			name:  "builder and args",
			style: sqlf.BindStyleQuestion,
			fragment: sqlf.F(
				"WHERE ?=?",
				sqlf.F("id"),
				nil,
			),
			want:     "WHERE id=?",
			wantArgs: []any{nil},
		},
		{
			name:  "build nil column",
			style: sqlf.BindStyleDollar,
			fragment: sqlf.F(
				"WHERE ?=?",
				(*sqlf.Fragment)(nil),
				nil,
			),
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name:  "args merging",
			style: sqlf.BindStyleDollar,
			fragment: sqlf.F(
				"WHERE foo=? AND bar IN (?)",
				1,
				sqlf.Join(",", 1, 2, 3),
			),
			want:     "WHERE foo=$1 AND bar IN ($1,$2,$3)",
			wantArgs: []any{1, 2, 3},
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.PrefixSuffix("SELECT", "FOR UPDATE", nil),
			want:     "",
			wantArgs: nil,
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.PrefixSuffix("SELECT", "FOR UPDATE", sqlf.F("foo")),
			want:     "SELECT foo FOR UPDATE",
			wantArgs: nil,
		},
		{
			name:  "ref fragment twice",
			style: sqlf.BindStyleDollar,
			fragment: sqlf.F(
				"$1, $1",
				sqlf.F("$1, $2, $1", 1, 2),
			),
			want:     "$1, $2, $1, $1, $2, $1",
			wantArgs: []any{1, 2},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			ctx := sqlf.NewContext(tc.style)
			got, err := tc.fragment.Build(ctx)
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
			args := ctx.Args()
			if !reflect.DeepEqual(args, tc.wantArgs) {
				t.Errorf("got %v, want %v", args, tc.wantArgs)
			}
		})
	}
}
