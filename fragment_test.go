package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/sqlb"
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

func TestBuildFragment(t *testing.T) {
	t.Parallel()
	var table, alias sqlb.Table = "table", "t"
	testCases := []struct {
		name     string
		style    syntax.BindVarStyle
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
			style: syntax.Question,
			fragment: sqlf.F(
				"WHERE ?=?",
				alias.Column("id"),
				nil,
			),
			want:     "WHERE t.id=?",
			wantArgs: []any{nil},
		},
		{
			name:  "build nil column",
			style: syntax.Dollar,
			fragment: sqlf.F(
				"WHERE ?=?",
				(*sqlf.Fragment)(nil),
				nil,
			),
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name:  "build complex fragment",
			style: syntax.Dollar,
			fragment: sqlf.F(
				"WITH t AS (?) SELECT ?,?,? FROM ? AS ?",
				sqlf.F(
					"SELECT * FROM ? AS ? WHERE ? > ?",
					table, alias, alias.Column("id"), 1,
				),
				alias.Column("id"),
				sqlf.F("?.id=?", alias, 2),
				"foo", table, alias,
			),
			want:     "WITH t AS (SELECT * FROM table AS t WHERE t.id > $1) SELECT t.id,t.id=$2,$3 FROM table AS t",
			wantArgs: []any{1, 2, "foo"},
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.F("").WithPrefix("SELECT").WithSuffix("FOR UPDATE"),
			want:     "",
			wantArgs: nil,
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.F("foo").WithPrefix("SELECT").WithSuffix("FOR UPDATE"),
			want:     "SELECT foo FOR UPDATE",
			wantArgs: nil,
		},
		{
			name:  "ref fragment twice",
			style: syntax.Dollar,
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
