package sqlf

import "fmt"

// Fragment is the builder for a part of or even a full query, it allows you
// to write and combine fragments with freedom.
type Fragment struct {
	Raw  string // Raw string support bind vars (?, $1)
	Args []any  // Args can be referenced by the Raw, for example: ?, $1. An arg can be either a sql arg or fragment builder.
}

// F creates a new Fragment.
// `raw` has exactly the same bind var (`?` / `$n`)
// syntax as `database/sql`, but more than that,
// it allows you to bind other fragment builders.
func F(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// WithArgs sets the args of f.
func (f *Fragment) WithArgs(args ...any) *Fragment {
	f.Args = args
	return f
}

// AppendArgs appends args to f.
func (f *Fragment) AppendArgs(args ...any) *Fragment {
	f.Args = append(f.Args, args...)
	return f
}

// SetArg sets the args at the given index (0-based, -1 means last).
func (f *Fragment) SetArg(index int, args any) *Fragment {
	if index < 0 {
		index = len(f.Args) + index
	}
	if index < 0 || index >= len(f.Args) {
		panic(fmt.Errorf("index %d out of range [%d,%d)", index, -len(f.Args), len(f.Args)))
	}
	f.Args[index] = args
	return f
}
