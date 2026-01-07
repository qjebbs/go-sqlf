package sqlf_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v4"
)

func TestBuildFragmentFn(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		style    sqlf.BindStyle
		builder  sqlf.Builder
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "join",
			style:    sqlf.BindStyleQuestion,
			builder:  sqlf.JoinArgs(",", 1, 2),
			want:     "?,?",
			wantArgs: []any{1, 2},
		},
		{
			name:    "prefix",
			style:   sqlf.BindStyleDollar,
			builder: sqlf.Prefix("WHERE", sqlf.F("1=1")),
			want:    "WHERE 1=1",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			ctx := sqlf.NewContext(context.Background(), tc.style)
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
