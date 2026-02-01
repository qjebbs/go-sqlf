package sqlf_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlf/v4/dialect"
)

func TestWithContextFunc(t *testing.T) {
	parent := sqlf.NewContext(context.Background(), dialect.PostgreSQL{})
	child := sqlf.ContextWithValue(parent, "k", "v")
	typeParent := reflect.TypeOf(parent)
	typeChild := reflect.TypeOf(child)

	if typeParent != typeChild {
		t.Fatalf("expected child context to have same type as parent: got %v, want %v", typeChild, typeParent)
	}
	if value := child.Value("k"); value != "v" {
		t.Fatalf("expected context value to be 'v': got %v", value)
	}
}

func TestNewContextValues(t *testing.T) {
	want := dialect.PostgreSQL{}
	ctx := sqlf.NewContext(context.Background(), want)
	got := ctx.BaseDialect()
	if got != want {
		t.Fatal("BaseDialect returned wrong value")
	}

	ctx.CommitArg(1)
	wantArgs := []any{1}
	gotArgs := ctx.Args()
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("Args returned wrong value: got %v, want %v", gotArgs, wantArgs)
	}
}
