package syntax

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	newExpr := func(line, col uint) expr {
		return expr{node{Pos{line, col}}}
	}
	testCases := []struct {
		raw     string
		want    []Expr
		wantErr bool
	}{
		{
			raw:     "?1",
			wantErr: true,
		},
		{
			raw:     "$",
			wantErr: true,
		},
		{
			raw:     "$1,?",
			wantErr: true,
		},
		{
			raw: "?,?,?",
			want: []Expr{
				&BindVarExpr{Type: Question, Index: 1, expr: newExpr(1, 1)},
				&PlainExpr{Text: ",", expr: newExpr(1, 2)},
				&BindVarExpr{Type: Question, Index: 2, expr: newExpr(1, 3)},
				&PlainExpr{Text: ",", expr: newExpr(1, 4)},
				&BindVarExpr{Type: Question, Index: 3, expr: newExpr(1, 5)},
			},
		},
		{
			raw: "$1'#c11#t111#fragment1111'",
			want: []Expr{
				&BindVarExpr{Type: Dollar, Index: 1, expr: newExpr(1, 1)},
				&PlainExpr{Text: "'#c11#t111#fragment1111'", expr: newExpr(1, 3)},
			},
		},
		{
			raw: "#join('#c=#argDollar', ',')",
			want: []Expr{
				&FuncCallExpr{
					Name: "join",
					Args: []any{"#c=#argDollar", ","},
					expr: newExpr(1, 1),
				},
			},
		},
		{
			raw: "#c1#t1#fragment1",
			want: []Expr{
				&FuncCallExpr{Name: "c", Args: []any{float64(1)}, expr: newExpr(1, 1)},
				&FuncCallExpr{Name: "t", Args: []any{float64(1)}, expr: newExpr(1, 4)},
				&FuncCallExpr{Name: "fragment", Args: []any{float64(1)}, expr: newExpr(1, 7)},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := Parse(tc.raw)
			if !tc.wantErr && err != nil {
				t.Fatal(err)
			}
			if !tc.wantErr && !reflect.DeepEqual(got.ExprList, tc.want) {
				for _, tk := range got.ExprList {
					t.Logf("%#v", tk)
				}
				t.Fatal("failed")
			}
		})
	}
}
