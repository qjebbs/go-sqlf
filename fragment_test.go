package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/sqlb"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func TestBuildFragment(t *testing.T) {
	t.Parallel()
	var table, alias sqlb.Table = "table", "t"
	testCases := []struct {
		name     string
		style    syntax.BindVarStyle
		fragment *sqlf.Fragment
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "build nil fragment",
			fragment: nil,
			want:     "",
			wantArgs: []any{},
		},
		{
			name:  "#join",
			style: syntax.Question,
			fragment: sqlf.Fa(
				"?,#join('#arg',',')",
				1, 2,
			),
			want:     "?,?,?",
			wantArgs: []any{1, 1, 2},
		},
		{
			name:  "#join range",
			style: syntax.Dollar,
			fragment: sqlf.Fa(
				"$1,#join('#arg',',', 2)",
				1, 2, 3, 4,
			),
			want:     "$1,$2,$3,$4",
			wantArgs: []any{1, 2, 3, 4},
		},
		{
			name:  "#join mixed function and call",
			style: syntax.Dollar,
			fragment: sqlf.F("#join('#f1#arg',',')").
				WithFragments(sqlf.Fa("p")).
				WithArgs(1, 2),
			want:     "p$1,p$2",
			wantArgs: []any{1, 2},
		},
		{
			name: "#f",
			fragment: sqlf.Ff("WHERE 1=1 #f1").
				WithFragments(sqlf.F("")),
			want:     "WHERE 1=1",
			wantArgs: []any{},
		},
		{
			name:  "#f and args",
			style: syntax.Question,
			fragment: sqlf.F("WHERE #f1=?").
				WithFragments(alias.Column("id")).
				WithArgs(nil),
			want:     "WHERE t.id=?",
			wantArgs: []any{nil},
		},
		{
			name:  "build nil column",
			style: syntax.Dollar,
			fragment: sqlf.F("WHERE #f1=$1").
				WithFragments((*sqlf.Fragment)(nil)).
				WithArgs(nil),
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name:  "build complex fragment",
			style: syntax.Dollar,
			fragment: sqlf.F("WITH t AS (#f1) SELECT #f2,#f3,$1 FROM #f4 AS #f5").
				WithArgs("foo").
				WithFragments(
					sqlf.F("SELECT * FROM #f1 AS #f2 WHERE #f3 > $1").
						WithArgs(1).
						WithFragments(table, alias, alias.Column("id")),
					alias.Column("id"),
					sqlf.F("#f1.id=$1").WithFragments(alias).WithArgs(2),
					table, alias,
				),
			want:     "WITH t AS (SELECT * FROM table AS t WHERE t.id > $1) SELECT t.id,t.id=$2,$3 FROM table AS t",
			wantArgs: []any{1, 2, "foo"},
		},
		{
			name:  "build complex fragment 2",
			style: syntax.Dollar,
			fragment: sqlf.F("SELECT #join('#f', ', ', 3) FROM #f1 AS #f2").
				WithFragments(
					table, alias,
					alias.Column("id"),
					sqlf.F("#f1.id=$1").WithFragments(alias).WithArgs(1),
					alias.Column("name"),
				),
			want:     "SELECT t.id, t.id=$1, t.name FROM table AS t",
			wantArgs: []any{1},
		},
		{
			name: "prefix and suffix",
			fragment: sqlf.F("#f1").WithPrefix("SELECT").WithSuffix("FOR UPDATE").
				WithFragments(sqlf.F("")),
			want:     "",
			wantArgs: []any{},
		},
		{
			name: "prefix and suffix",
			fragment: sqlf.F("#f1").WithPrefix("SELECT").WithSuffix("FOR UPDATE").
				WithFragments(sqlf.F("foo")),
			want:     "SELECT foo FOR UPDATE",
			wantArgs: []any{},
		},
		{
			name:  "ref fragment twice",
			style: syntax.Dollar,
			fragment: sqlf.F("#f1, #f1").
				WithFragments(
					sqlf.F("#join('#arg', ', '), ?").WithArgs(1, 2),
				),
			want:     "$1, $2, $1, $1, $2, $1",
			wantArgs: []any{1, 2},
		},
		{
			name:  "arg and fragment",
			style: syntax.Question,
			fragment: sqlf.Fa("? #f1", 1).
				WithFragments(
					sqlf.Fa("$1", 2),
				),
			want:     "? ?",
			wantArgs: []any{1, 2},
		},
		{
			name:     "mixed bindvar style",
			fragment: sqlf.Fa("?, $1", nil),
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			ctx := sqlf.NewContext(tc.style)
			got, err := tc.fragment.BuildFragment(ctx)
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
