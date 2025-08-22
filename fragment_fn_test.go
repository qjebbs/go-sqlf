package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func TestBuildFragmentFn(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		style    syntax.BindVarStyle
		builder  sqlf.Builder
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "join",
			style:    syntax.Question,
			builder:  sqlf.Join(",", 1, 2),
			want:     "?,?",
			wantArgs: []any{1, 2},
		},
		{
			name:  "join 2",
			style: syntax.Question,
			builder: sqlf.Fn(func(ctx *sqlf.Context) (query string, err error) {
				args := []any{1, 2, 3, 4}
				return sqlf.F(
					"WHERE foo=? AND bar IN (?)",
					args[0],
					sqlf.Join(",", args[1:]...),
				).BuildFragment(ctx)
			}),
			want:     "WHERE foo=? AND bar IN (?,?,?)",
			wantArgs: []any{1, 2, 3, 4},
		},
		{
			name:  "args merging",
			style: syntax.Dollar,
			builder: sqlf.Fn(func(ctx *sqlf.Context) (query string, err error) {
				args := []any{1, 1, 2, 3}
				return sqlf.F(
					"WHERE foo=? AND bar IN (?)",
					args[0],
					sqlf.Join(",", args[1:]...),
				).BuildFragment(ctx)
			}),
			want:     "WHERE foo=$1 AND bar IN ($1,$2,$3)",
			wantArgs: []any{1, 2, 3},
		},
		{
			name:  "build complex fragment",
			style: syntax.Dollar,
			builder: sqlf.Fn(func(ctx *sqlf.Context) (query string, err error) {
				t := sqlb.NewTableAliased("foo", "f")
				return sqlf.F(
					"SELECT $1,$1=$2 FROM $3 AS $4",
					t.Column("id"), 1,
					t.Name, t.Alias,
				).BuildFragment(ctx)
			}),
			want:     "SELECT f.id,f.id=$1 FROM foo AS f",
			wantArgs: []any{1},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			ctx := sqlf.NewContext(tc.style)
			got, err := tc.builder.BuildFragment(ctx)
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
