package sqlf

import (
	"strings"

	"github.com/qjebbs/go-sqlf/v2/syntax"
)

// Join creates a new fragment builder that joins the given arguments with the specified separator.
//
// An arg could be either a sql arg or fragment builder.
func Join(sep string, args ...any) Builder {
	return Fn(func(ctx *Context) (string, error) {
		if len(args) == 0 {
			return "", nil
		}
		var sb strings.Builder
		props := newProperties(args...)
		for i, p := range props {
			if i > 0 {
				sb.WriteString(sep)
			}
			r, err := p.BuildFragment(ctx)
			if err != nil {
				return "", err
			}
			sb.WriteString(r)
		}
		return sb.String(), nil
	})
}

// Fn creates a new fragment builder with the fn function.
func Fn(fn func(ctx *Context) (query string, err error)) Builder {
	return &builder{
		fn: fn,
	}
}

var _ Builder = (*builder)(nil)
var _ QueryBuilder = (*builder)(nil)

type builder struct {
	fn func(ctx *Context) (query string, err error)
}

func (b *builder) BuildFragment(ctx *Context) (string, error) {
	return b.fn(ctx)
}

// BuildQuery builds the fragment as full query.
func (b *builder) BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error) {
	return _buildBuilder(b, bindVarStyle)
}
