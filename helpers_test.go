package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/syntax"
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
			name:    "prefix",
			style:   syntax.Dollar,
			builder: sqlf.Prefix("WHERE", sqlf.F("1=1")),
			want:    "WHERE 1=1",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			ctx := sqlf.NewContext(tc.style)
			got, err := tc.builder.Build(ctx)
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
