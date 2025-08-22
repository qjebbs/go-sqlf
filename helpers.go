package sqlf

import (
	"strings"
)

// Join creates a new fragment builder that joins the given arguments with the specified separator.
//
// An arg could be either a sql arg or fragment builder.
func Join(sep string, args ...any) Builder {
	return fn(func(ctx *Context) (string, error) {
		if len(args) == 0 {
			return "", nil
		}
		var sb strings.Builder
		props := newProperties(args...)
		for i, p := range props {
			if i > 0 {
				sb.WriteString(sep)
			}
			r, err := p.Build(ctx)
			if err != nil {
				return "", err
			}
			sb.WriteString(r)
		}
		return sb.String(), nil
	})
}

// Prefix creates a new fragment builder that prefixes the given builder if it's built not empty.
func Prefix(prefix string, b Builder) Builder {
	return fn(func(ctx *Context) (query string, err error) {
		query, err = b.Build(ctx)
		if err != nil {
			return "", err
		}
		if query == "" {
			return "", nil
		}
		return prefix + " " + query, nil
	})
}

// fn creates a new fragment builder with the fn function.
func fn(fn func(ctx *Context) (query string, err error)) Builder {
	return &builder{
		fn: fn,
	}
}

var _ Builder = (*builder)(nil)

type builder struct {
	fn func(ctx *Context) (query string, err error)
}

func (b *builder) Build(ctx *Context) (string, error) {
	return b.fn(ctx)
}
