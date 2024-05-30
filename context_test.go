package sqlf

import (
	"strings"
	"testing"

	"github.com/qjebbs/go-sqlf/v2/syntax"
)

func TestContextWithFragment(t *testing.T) {
	t.Parallel()
	fragment := Ff(
		"L1,#f1", Ff(
			"L2,#f1", Ff(
				"L3,#f1", Ff("L4,#parents()"),
			),
		),
	)
	ctx := NewContext(syntax.Dollar)
	ctx, err := ContextWithFuncs(ctx, FuncMap{
		"parents": func(ctx *Context) (string, error) {
			parents := make([]string, 0)
			for c := ctx.parent; c != nil; c = c.parent {
				fc := c.Fragment()
				if fc == nil {
					continue
				}
				parents = append(parents, strings.SplitN(fc.Fragment.Raw, ",", 2)[0])
			}
			return strings.Join(parents, ","), nil
		},
	})
	if err != nil {
		t.Fatalf("WithFuncs failed: %v", err)
	}
	want := "L1,L2,L3,L4,L3,L2,L1"
	got, err := fragment.BuildFragment(ctx)
	if err != nil {
		t.Fatalf("BuildFragment failed: %v", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
