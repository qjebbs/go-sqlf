package sqlf

import (
	"strings"
)

// Join creates a new fragment builder that joins the given builders with the specified separator.
//
// An arg could be either an ordinary arg or a Builder.
func Join(builders []Builder, sep string) Builder {
	return Func(func(ctx Context) (string, error) {
		if len(builders) == 0 {
			return "", nil
		}
		var sb strings.Builder
		for i, p := range builders {
			r, err := p.BuildTo(ctx)
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

// JoinArgs creates a new fragment builder that joins args with the specified separator.
func JoinArgs[T any](args []T, sep string) Builder {
	return Func(func(ctx Context) (string, error) {
		if len(args) == 0 {
			return "", nil
		}
		props := make([]Builder, 0, len(args))
		for _, a := range args {
			props = append(props, newArgProperty(a))
		}
		return Join(props, sep).BuildTo(ctx)
	})
}

// JoinMixed creates a new fragment builder that joins mixed Builders and args with the specified separator.
func JoinMixed(args []any, sep string) Builder {
	return Func(func(ctx Context) (string, error) {
		if len(args) == 0 {
			return "", nil
		}
		props := make([]Builder, 0, len(args))
		for _, a := range args {
			if b, ok := a.(Builder); ok {
				props = append(props, b)
				continue
			}
			props = append(props, newArgProperty(a))
		}
		return Join(props, sep).BuildTo(ctx)
	})
}

// Prefix creates a new fragment builder that prefixes the given builder if it's built not empty.
func Prefix(prefix string, b Builder) Builder {
	return Func(func(ctx Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
		query, err = b.BuildTo(ctx)
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
	return Func(func(ctx Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
		query, err = b.BuildTo(ctx)
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
	return Func(func(ctx Context) (query string, err error) {
		if b == nil {
			return "", nil
		}
		query, err = b.BuildTo(ctx)
		if err != nil {
			return "", err
		}
		if query == "" {
			return "", nil
		}
		return prefix + " " + query + " " + suffix, nil
	})
}

// Identifier creates a new fragment builder that quotes the given identifier using current dialect.
func Identifier(name string) Builder {
	return Func(func(ctx Context) (string, error) {
		dialect := ctx.BaseDialect()
		quoted := dialect.QuoteIdentifier(name)
		return quoted, nil
	})
}

// Func is a helper to create fragment builder with function.
func Func(fn func(ctx Context) (query string, err error)) Builder {
	return &builder{
		fn: fn,
	}
}

var _ Builder = (*builder)(nil)

type builder struct {
	fn func(ctx Context) (query string, err error)
}

func (b *builder) BuildTo(ctx Context) (string, error) {
	if b == nil || b.fn == nil {
		return "", nil
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return b.fn(ctx)
}
