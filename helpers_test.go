package sqlf_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/arg"
	"github.com/qjebbs/go-sqlf/v4/dialect"
)

func TestBuildFragmentFn(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		builder  sqlf.Builder
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "join",
			builder:  sqlf.JoinArgs(",", 1, 2),
			want:     "?,?",
			wantArgs: []any{1, 2},
		},
		{
			name:    "prefix",
			builder: sqlf.Prefix("WHERE", sqlf.F("1=1")),
			want:    "WHERE 1=1",
		},
	}
	ctx := sqlf.NewContext(context.Background(), dialect.SQLite{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			got, args, err := sqlf.Build(ctx, tc.builder)
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

func TestBuildIdentifiers(t *testing.T) {
	testCases := []struct {
		name    string
		dialect dialect.Dialect
		ident   string
		want    string
	}{
		{
			name:    "PostgreSQL",
			dialect: dialect.PostgreSQL{},
			ident:   `user"name`,
			want:    `"user""name"`,
		},
		{
			name:    "SQLServer",
			dialect: dialect.SQLServer{},
			ident:   `user [name]`,
			want:    `[user [name]]]`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := sqlf.NewContext(context.Background(), tc.dialect)
			builder := sqlf.Identifier(tc.ident)
			got, err := builder.BuildTo(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	f := sqlf.JoinArgs(", ", 1, 2, 3)
	ctx := sqlf.NewContext(context.Background(), dialect.SQLite{
		BindVarStyle: arg.StyleDollarNumbered,
	})
	got, _, err := sqlf.Build(ctx, f)
	if err != nil {
		t.Fatal(err)
	}
	want := "$1, $2, $3"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if n := len(ctx.Args()); n != 0 {
		t.Errorf("expected no args committed to the store, got %d", n)
	}
}
