package sqlf

import (
	"context"
	"testing"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/arg"
)

func TestContextValues(t *testing.T) {
	testCases := []struct {
		name string
		fn   func(context.Context) *Context
	}{
		{
			name: "NewContext",
			fn: func(ctx context.Context) *Context {
				return NewContext(ctx, dialect.PostgreSQL{})
			},
		},
		{
			name: "ContextWithDialect",
			fn: func(ctx context.Context) *Context {
				return ContextWithDialect(ctx, dialect.PostgreSQL{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.fn(context.Background())
			assertCongextValuest(t, ctx)
		})
	}
}

func assertCongextValuest(t *testing.T, ctx *Context) {
	t.Helper()
	value := ctx.Value(dialectKey{})
	if value == nil {
		t.Fatal("Dialect not found in context")
	}
	_, ok := value.(dialect.Dialect)
	if !ok {
		t.Fatal("Dialect not found in context")
	}
	value = ctx.Value(argStoreKey{})
	if value == nil {
		t.Fatal("ArgStore not found in context")
	}
	_, ok = value.(arg.Store)
	if !ok {
		t.Fatal("ArgStore has wrong type")
	}
}
