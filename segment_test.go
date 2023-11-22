package sqls_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqls"
)

func TestBuildSegment(t *testing.T) {
	t.Parallel()
	var table, alias sqls.Table = "table", "t"
	testCases := []struct {
		name     string
		segment  *sqls.Segment
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "build nil segment",
			segment:  nil,
			want:     "",
			wantArgs: []any{},
		},
		{
			name: "#join",
			segment: &sqls.Segment{
				Raw:  "#join('#?',','),#?(1),#?(2)",
				Args: []any{1, 2},
			},
			want:     "?,?,?,?",
			wantArgs: []any{1, 2, 1, 2},
		},
		{
			name: "#join mixed function and call",
			segment: &sqls.Segment{
				Raw:      "#join('#s1#?',',')",
				Args:     []any{1, 2},
				Segments: []*sqls.Segment{{Raw: "s1"}},
			},
			want:     "s1?,s1?",
			wantArgs: []any{1, 2},
		},
		{
			name: "#segment",
			segment: &sqls.Segment{
				Raw:      "WHERE 1=1 #s1",
				Segments: []*sqls.Segment{nil},
			},
			want:     "WHERE 1=1",
			wantArgs: []any{},
		},
		{
			name: "#column and args",
			segment: &sqls.Segment{
				Raw:     "WHERE #c1=?",
				Columns: alias.Columns("id"),
				Args:    []any{nil},
			},
			want:     "WHERE t.id=?",
			wantArgs: []any{nil},
		},
		{
			name: "build nil column",
			segment: &sqls.Segment{
				Raw:     "WHERE #c1=$1",
				Columns: []*sqls.TableColumn{nil},
				Args:    []any{nil},
			},
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name: "build column without args",
			segment: &sqls.Segment{
				Raw:     "#c1>1",
				Columns: alias.Columns("id"),
				Args:    nil,
			},
			want:     "t.id>1",
			wantArgs: []any{},
		},
		{
			name: "build column with args",
			segment: &sqls.Segment{
				Raw:     "#c2 IS NULL AND #c1>$1",
				Columns: alias.Columns("id", "deleted"),
				Args:    []any{1},
			},
			want:     "t.deleted IS NULL AND t.id>$1",
			wantArgs: []any{1},
		},
		{
			name: "build column with args 2",
			segment: &sqls.Segment{
				Raw:     "#c1>$1",
				Columns: alias.Columns("id"),
				Args:    []any{1},
			},
			want:     "t.id>$1",
			wantArgs: []any{1},
		},
		{
			name: "build column with unusual args order",
			segment: &sqls.Segment{
				Raw:     "#c1 IN ($2,$1)",
				Columns: alias.Columns("id"),
				Args:    []any{1, 2},
			},
			want:     "t.id IN ($1,$2)",
			wantArgs: []any{2, 1},
		},
		{
			name: "build column expression with args",
			segment: &sqls.Segment{
				Raw: "#c1",
				Columns: []*sqls.TableColumn{
					alias.Expression("#t1.id=$1", 1),
				},
			},
			want:     "t.id=$1",
			wantArgs: []any{1},
		},
		{
			name: "build column expression with args, and args",
			segment: &sqls.Segment{
				Raw: "#c1 > $1",
				Columns: []*sqls.TableColumn{
					alias.Expression("#t1.id - $1", 1),
				},
				Args: []any{2},
			},
			want:     "t.id - $1 > $2",
			wantArgs: []any{1, 2},
		},
		{
			name: "build complex segment",
			segment: &sqls.Segment{
				Raw: "WITH t AS (#s1) SELECT #c1,#c2,$1 FROM #t1 AS #t2 ",
				Segments: []*sqls.Segment{
					{
						Raw:     "SELECT * FROM #t1 AS #t2 WHERE #c1 > $1",
						Columns: alias.Columns("id"),
						Tables:  []sqls.Table{table, alias},
						Args:    []any{1},
					},
				},
				Columns: []*sqls.TableColumn{
					alias.Column("id"),
					alias.Expression("#t1.id=$1", 2),
				},
				Tables: []sqls.Table{table, alias},
				Args:   []any{"foo"},
			},
			want:     "WITH t AS (SELECT * FROM table AS t WHERE t.id > $1) SELECT t.id,t.id=$2,$3 FROM table AS t",
			wantArgs: []any{1, 2, "foo"},
		},
		{
			name: "build complex segment 2",
			segment: &sqls.Segment{
				Raw: "SELECT #join('#c', ', ') FROM #t1 AS #t2 ",
				Columns: []*sqls.TableColumn{
					alias.Column("id"),
					alias.Expression("#t1.id=$1", 1),
					alias.Column("name"),
				},
				Tables: []sqls.Table{table, alias},
			},
			want:     "SELECT t.id, t.id=$1, t.name FROM table AS t",
			wantArgs: []any{1},
		},
		{
			name: "prefix and suffix",
			segment: &sqls.Segment{
				Raw:      "#s1",
				Segments: []*sqls.Segment{nil},
				Prefix:   "WHERE",
				Suffix:   "FOR UPDATE",
			},
			want:     "",
			wantArgs: []any{},
		},
		{
			name: "prefix and suffix deep",
			segment: &sqls.Segment{
				Raw: "#s1",
				Segments: []*sqls.Segment{
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
			name: "ref segment twice",
			segment: &sqls.Segment{
				Raw: "#s1, #s1",
				Segments: []*sqls.Segment{{
					Raw:  "#join('#?', ', '), ?",
					Args: []any{1, 2},
				}},
			},
			want:     "?, ?, ?, ?, ?, ?",
			wantArgs: []any{1, 2, 1, 1, 2, 1},
		},
		{
			name: "arg and segment",
			segment: &sqls.Segment{
				Raw: "? #s1",
				Segments: []*sqls.Segment{{
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
			segment: &sqls.Segment{
				Raw:  "?, $1",
				Args: []any{nil},
			},
			wantErr: true,
		},
		{
			name: "build builder",
			segment: &sqls.Segment{
				Raw: "id IN (#b1)",
				Builders: []sqls.Builder{
					&sqls.Segment{
						Raw:     "SELECT id FROM #t1 WHERE #c1 > $1",
						Tables:  []sqls.Table{table},
						Columns: alias.Expressions("id"),
						Args:    []any{1},
					},
				},
			},
			want:     "id IN (SELECT id FROM table WHERE id > $1)",
			wantArgs: []any{1},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			args := make([]any, 0)
			got, err := tc.segment.BuildContext(sqls.NewContext(&args))
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
