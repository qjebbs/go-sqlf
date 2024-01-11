package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf"
)

func TestBuildFragment(t *testing.T) {
	t.Parallel()
	var table, alias sqlf.Table = "table", "t"
	testCases := []struct {
		name       string
		fragment   *sqlf.Fragment
		globalArgs []any
		want       string
		wantArgs   []any
		wantErr    bool
	}{
		{
			name:     "build nil fragment",
			fragment: nil,
			want:     "",
			wantArgs: []any{},
		},
		{
			name: "#join",
			fragment: &sqlf.Fragment{
				Raw:  "#join('#?',','),#?(1),#?(2)",
				Args: []any{1, 2},
			},
			want:     "?,?,?,?",
			wantArgs: []any{1, 2, 1, 2},
		},
		{
			name: "#join range",
			fragment: &sqlf.Fragment{
				Raw:  "$1,#join('#$',',', 2)",
				Args: []any{1, 2, 3, 4},
			},
			want:     "$1,$2,$3,$4",
			wantArgs: []any{1, 2, 3, 4},
		},
		{
			name: "#join mixed function and call",
			fragment: &sqlf.Fragment{
				Raw:       "#join('#f1#?',',')",
				Args:      []any{1, 2},
				Fragments: []*sqlf.Fragment{{Raw: "s1"}},
			},
			want:     "s1?,s1?",
			wantArgs: []any{1, 2},
		},
		{
			name: "#fragment",
			fragment: &sqlf.Fragment{
				Raw:       "WHERE 1=1 #f1",
				Fragments: []*sqlf.Fragment{nil},
			},
			want:     "WHERE 1=1",
			wantArgs: []any{},
		},
		{
			name: "#column and args",
			fragment: &sqlf.Fragment{
				Raw:     "WHERE #c1=?",
				Columns: alias.Columns("id"),
				Args:    []any{nil},
			},
			want:     "WHERE t.id=?",
			wantArgs: []any{nil},
		},
		{
			name: "build nil column",
			fragment: &sqlf.Fragment{
				Raw:     "WHERE #c1=$1",
				Columns: []*sqlf.TableColumn{nil},
				Args:    []any{nil},
			},
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name: "build column without args",
			fragment: &sqlf.Fragment{
				Raw:     "#c1>1",
				Columns: alias.Columns("id"),
				Args:    nil,
			},
			want:     "t.id>1",
			wantArgs: []any{},
		},
		{
			name: "build column with args",
			fragment: &sqlf.Fragment{
				Raw:     "#c2 IS NULL AND #c1>$1",
				Columns: alias.Columns("id", "deleted"),
				Args:    []any{1},
			},
			want:     "t.deleted IS NULL AND t.id>$1",
			wantArgs: []any{1},
		},
		{
			name: "build column with args 2",
			fragment: &sqlf.Fragment{
				Raw:     "#c1>$1",
				Columns: alias.Columns("id"),
				Args:    []any{1},
			},
			want:     "t.id>$1",
			wantArgs: []any{1},
		},
		{
			name: "build column with unusual args order",
			fragment: &sqlf.Fragment{
				Raw:     "#c1 IN ($2,$1)",
				Columns: alias.Columns("id"),
				Args:    []any{1, 2},
			},
			want:     "t.id IN ($1,$2)",
			wantArgs: []any{2, 1},
		},
		{
			name: "build column expression with args",
			fragment: &sqlf.Fragment{
				Raw: "#c1",
				Columns: []*sqlf.TableColumn{
					alias.Expression("#t1.id=$1", 1),
				},
			},
			want:     "t.id=$1",
			wantArgs: []any{1},
		},
		{
			name: "build column expression with args, and args",
			fragment: &sqlf.Fragment{
				Raw: "#c1 > $1",
				Columns: []*sqlf.TableColumn{
					alias.Expression("#t1.id - $1", 1),
				},
				Args: []any{2},
			},
			want:     "t.id - $1 > $2",
			wantArgs: []any{1, 2},
		},
		{
			name: "build complex fragment",
			fragment: &sqlf.Fragment{
				Raw: "WITH t AS (#f1) SELECT #c1,#c2,$1 FROM #t1 AS #t2 ",
				Fragments: []*sqlf.Fragment{
					{
						Raw:     "SELECT * FROM #t1 AS #t2 WHERE #c1 > $1",
						Columns: alias.Columns("id"),
						Tables:  []sqlf.Table{table, alias},
						Args:    []any{1},
					},
				},
				Columns: []*sqlf.TableColumn{
					alias.Column("id"),
					alias.Expression("#t1.id=$1", 2),
				},
				Tables: []sqlf.Table{table, alias},
				Args:   []any{"foo"},
			},
			want:     "WITH t AS (SELECT * FROM table AS t WHERE t.id > $1) SELECT t.id,t.id=$2,$3 FROM table AS t",
			wantArgs: []any{1, 2, "foo"},
		},
		{
			name: "build complex fragment 2",
			fragment: &sqlf.Fragment{
				Raw: "SELECT #join('#c', ', ') FROM #t1 AS #t2 ",
				Columns: []*sqlf.TableColumn{
					alias.Column("id"),
					alias.Expression("#t1.id=$1", 1),
					alias.Column("name"),
				},
				Tables: []sqlf.Table{table, alias},
			},
			want:     "SELECT t.id, t.id=$1, t.name FROM table AS t",
			wantArgs: []any{1},
		},
		{
			name: "prefix and suffix",
			fragment: &sqlf.Fragment{
				Raw:       "#f1",
				Fragments: []*sqlf.Fragment{nil},
				Prefix:    "WHERE",
				Suffix:    "FOR UPDATE",
			},
			want:     "",
			wantArgs: []any{},
		},
		{
			name: "prefix and suffix deep",
			fragment: &sqlf.Fragment{
				Raw: "#f1",
				Fragments: []*sqlf.Fragment{
					{
						Raw:     "#c1=$1",
						Columns: alias.Columns("id"),
						Args:    []any{1},
					},
				},
				Prefix: "WHERE",
				Suffix: "FOR UPDATE",
			},
			want:     "WHERE t.id=$1 FOR UPDATE",
			wantArgs: []any{1},
		},
		{
			name: "ref fragment twice",
			fragment: &sqlf.Fragment{
				Raw: "#f1, #f1",
				Fragments: []*sqlf.Fragment{{
					Raw:  "#join('#?', ', '), ?",
					Args: []any{1, 2},
				}},
			},
			want:     "?, ?, ?, ?, ?, ?",
			wantArgs: []any{1, 2, 1, 1, 2, 1},
		},
		{
			name: "arg and fragment",
			fragment: &sqlf.Fragment{
				Raw: "? #f1",
				Fragments: []*sqlf.Fragment{{
					Raw:  "$1",
					Args: []any{2},
				}},
				Args: []any{1},
			},
			want:     "? ?",
			wantArgs: []any{1, 2},
		},
		{
			name: "mixed bindvar style",
			fragment: &sqlf.Fragment{
				Raw:  "?, $1",
				Args: []any{nil},
			},
			wantErr: true,
		},
		{
			name: "build builder",
			fragment: &sqlf.Fragment{
				Raw: "id IN (#b1)",
				Builders: []sqlf.Builder{
					&sqlf.Fragment{
						Raw:     "SELECT id FROM #t1 WHERE #c1 > $1",
						Tables:  []sqlf.Table{table},
						Columns: alias.Expressions("id"),
						Args:    []any{1},
					},
				},
			},
			want:     "id IN (SELECT id FROM table WHERE id > $1)",
			wantArgs: []any{1},
		},
		{
			name: "build with global args $",
			fragment: &sqlf.Fragment{
				Raw: "#join('#fragment',' ')",
				Fragments: []*sqlf.Fragment{
					{Raw: "#global$1"},
					{Raw: "#global$2"},
					{Raw: "#global$2"},
				},
			},
			globalArgs: []any{1, 2},
			want:       "$1 $2 $2",
			wantArgs:   []any{1, 2},
		},
		{
			name: "build with global args ?",
			fragment: &sqlf.Fragment{
				Raw: "#join('#fragment',' ')",
				Fragments: []*sqlf.Fragment{
					{Raw: "#global?1"},
					{Raw: "#global?2"},
					{Raw: "#global?2"},
				},
			},
			globalArgs: []any{1, 2},
			want:       "? ? ?",
			wantArgs:   []any{1, 2, 2},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			args := make([]any, 0)
			got, err := tc.fragment.BuildContext(sqlf.NewContext(&args).WithArgs(tc.globalArgs))
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
			if !reflect.DeepEqual(args, tc.wantArgs) {
				t.Errorf("got %v, want %v", args, tc.wantArgs)
			}
		})
	}
}
