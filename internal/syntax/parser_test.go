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
			raw:     "$",
			wantErr: true,
		},
		{
			raw:     "$1,?",
			wantErr: true,
		},
		{
			raw: "??",
			want: []Expr{
				&PlainExpr{Text: "?", expr: newExpr(1, 1)},
			},
		},
		{
			raw: "$$",
			want: []Expr{
				&PlainExpr{Text: "$", expr: newExpr(1, 1)},
			},
		},
		{
			raw: "?1",
			want: []Expr{
				&BindVarExpr{typ: BindStyleQuestion, Index: 1, expr: newExpr(1, 1)},
				&PlainExpr{Text: "1", expr: newExpr(1, 2)},
			},
		},
		{
			raw: "?,?,?",
			want: []Expr{
				&BindVarExpr{typ: BindStyleQuestion, Index: 1, expr: newExpr(1, 1)},
				&PlainExpr{Text: ",", expr: newExpr(1, 2)},
				&BindVarExpr{typ: BindStyleQuestion, Index: 2, expr: newExpr(1, 3)},
				&PlainExpr{Text: ",", expr: newExpr(1, 4)},
				&BindVarExpr{typ: BindStyleQuestion, Index: 3, expr: newExpr(1, 5)},
			},
		},
		{
			raw: "$1'?,?,$1'",
			want: []Expr{
				&BindVarExpr{typ: BindStyleDollarNumbered, Index: 1, expr: newExpr(1, 1)},
				&PlainExpr{Text: "'?,?,$1'", expr: newExpr(1, 3)},
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
