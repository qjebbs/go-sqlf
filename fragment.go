package sqlf

import "fmt"

// Fragment is a builder for a part of, or even an entire, SQL query.
//
// 'Raw' uses the same bind variable syntax (? / $1) as database/sql.
// In addition, it supports binding other fragment builders.
type Fragment struct {
	Raw  string // Raw string support bind vars (?, $1)
	Args []any  // Args that can be referenced by the Raw. An arg can be either an ordinary arg or a Builder.
}

// F creates a new Fragment.
// 'raw' uses the same bind variable syntax (? / $1) as database/sql.
// Additionally, it allows you to bind other fragment builders.
func F(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// WithArgs replaces the args of f.
func (f *Fragment) WithArgs(args ...any) *Fragment {
	f.Args = args
	return f
}

// AppendArgs appends args to f.
func (f *Fragment) AppendArgs(args ...any) *Fragment {
	f.Args = append(f.Args, args...)
	return f
}

// SetArg sets the arg at the given index (0-based, -1 means last).
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
