package syntax_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/qjebbs/go-sqlf/syntax"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		raw     string
		want    []syntax.Expr
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
			want: []syntax.Expr{
				&syntax.BindVarExpr{Type: syntax.Question, Index: 1},
				&syntax.PlainExpr{Text: ","},
				&syntax.BindVarExpr{Type: syntax.Question, Index: 2},
				&syntax.PlainExpr{Text: ","},
				&syntax.BindVarExpr{Type: syntax.Question, Index: 3},
			},
		},
		{
			raw: "$1'#c11#t111#fragment1111'",
			want: []syntax.Expr{
				&syntax.BindVarExpr{Type: syntax.Dollar, Index: 1},
				&syntax.PlainExpr{Text: "'#c11#t111#fragment1111'"},
			},
		},
		{
			raw: "#join('#c=#argDollar', ',')",
			want: []syntax.Expr{
				&syntax.FuncCallExpr{
					Name: "join",
					Args: []any{"#c=#argDollar", ","},
				},
			},
		},
		{
			raw: "#c1#t1#fragment1",
			want: []syntax.Expr{
				&syntax.FuncCallExpr{Name: "c", Args: []any{float64(1)}},
				&syntax.FuncCallExpr{Name: "t", Args: []any{float64(1)}},
				&syntax.FuncCallExpr{Name: "fragment", Args: []any{float64(1)}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := syntax.Parse(tc.raw)
			if !tc.wantErr && err != nil {
				t.Fatal(err)
			}
			if !tc.wantErr && !cmp.Equal(
				got.ExprList, tc.want,
				cmpopts.IgnoreUnexported(
					syntax.PlainExpr{},
					syntax.BindVarExpr{},
					syntax.FuncExpr{},
					syntax.FuncCallExpr{},
				),
			) {
				for _, tk := range got.ExprList {
					t.Logf("%#v", tk)
				}
				t.Fatal("failed")
			}
		})
	}
}
