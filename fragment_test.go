package sqlf_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/dialect"
)

func TestBuildFragment(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		fragment sqlf.Builder
		want     string
		wantArgs []any
		wantErr  bool
	}{
		{
			name:     "escaping",
			fragment: sqlf.F("WHERE foo -> ? ?? ?", "bar", "baz"),
			want:     "WHERE foo -> $1 ? $2",
			wantArgs: []any{"bar", "baz"},
		},
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
			name: "builder and args",
			fragment: sqlf.F(
				"WHERE ?=?",
				sqlf.F("id"),
				nil,
			),
			want:     "WHERE id=$1",
			wantArgs: []any{nil},
		},
		{
			name: "build nil column",
			fragment: sqlf.F(
				"WHERE ?=?",
				(*sqlf.Fragment)(nil),
				nil,
			),
			want:     "WHERE =$1",
			wantArgs: []any{nil},
		},
		{
			name: "args merging",
			fragment: sqlf.F(
				"WHERE foo=? AND bar IN (?)",
				1,
				sqlf.JoinArgs(",", 1, 2, 3),
			),
			want:     "WHERE foo=$1 AND bar IN ($1,$2,$3)",
			wantArgs: []any{1, 2, 3},
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.PrefixSuffix("SELECT", "FOR UPDATE", nil),
			want:     "",
			wantArgs: nil,
		},
		{
			name:     "prefix and suffix",
			fragment: sqlf.PrefixSuffix("SELECT", "FOR UPDATE", sqlf.F("foo")),
			want:     "SELECT foo FOR UPDATE",
			wantArgs: nil,
		},
		{
			name: "ref fragment twice",
			fragment: sqlf.F(
				"$1, $1",
				sqlf.F("$1, $2, $1", 1, 2),
			),
			want:     "$1, $2, $1, $1, $2, $1",
			wantArgs: []any{1, 2},
		},
	}
	ctx := sqlf.NewContext(context.Background(), dialect.PostgreSQL{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			got, args, err := sqlf.Build(ctx, tc.fragment)
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
func TestBuildWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	buildCtx := sqlf.NewContext(ctx, dialect.PostgreSQL{})
	_, err := sqlf.F(`SELECT 1`).BuildTo(buildCtx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got err %v, want context.Canceled", err)
	}
	_, err = sqlf.JoinArgs(",", 1, 2, 3).BuildTo(buildCtx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got err %v, want context.Canceled", err)
	}
}
