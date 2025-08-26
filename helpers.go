package sqlf

import (
	"strings"
)

// Join creates a new fragment builder that joins the given args with the specified separator.
//
// An arg could be either an ordinary arg or a Builder.
func Join(sep string, args ...any) Builder {
	return Func(func(ctx *Context) (string, error) {
		if len(args) == 0 {
			return "", nil
		}
		var sb strings.Builder
		props := newProperties(args...)
		for i, p := range props {
			r, err := p.Build(ctx)
			if err != nil {
				return "", err
			}
			r = strings.TrimSpace(r)
			if r == "" {
				continue
			}
			if i > 0 {
				sb.WriteString(sep)
			}
			sb.WriteString(r)
		}
		return sb.String(), nil
	})
}

// Prefix creates a new fragment builder that prefixes the given builder if it's built not empty.
func Prefix(prefix string, b Builder) Builder {
	return Func(func(ctx *Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
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

// Suffix creates a new fragment builder that suffixes the given builder if it's built not empty.
func Suffix(suffix string, b Builder) Builder {
	return Func(func(ctx *Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
		query, err = b.Build(ctx)
		if err != nil {
			return "", err
		}
		if query == "" {
			return "", nil
		}
		return query + " " + suffix, nil
	})
}

// PrefixSuffix creates a new fragment builder that prefixes and suffixes the given builder if it's built not empty.
func PrefixSuffix(prefix, suffix string, b Builder) Builder {
	return Func(func(ctx *Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
		query, err = b.Build(ctx)
		if err != nil {
			return "", err
		}
		if query == "" {
			return "", nil
		}
		return prefix + " " + query + " " + suffix, nil
	})
}

// Func is a helper to create fragment builder with function.
func Func(fn func(ctx *Context) (query string, err error)) Builder {
	return &builder{
		fn: fn,
	}
}

var _ Builder = (*builder)(nil)

type builder struct {
	fn func(ctx *Context) (query string, err error)
}

func (b *builder) Build(ctx *Context) (string, error) {
	if b == nil || b.fn == nil {
		return "", nil
	}
	return b.fn(ctx)
}
