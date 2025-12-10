package sqlf

import "fmt"

// Fragment is the builder for a part of, or even an entire, SQL query.
//
// To create a Fragment, use the F() function.
type Fragment struct {
	raw          string // Raw string support bind vars (?, $1)
	args         []any  // Args that can be referenced by the Raw. An arg can be either an ordinary arg or a Builder.
	noUsageCheck bool   // If true, skip checking whether all bind vars are used.
}

// F creates a new Fragment.
// 'raw' uses the same bind variable syntax (? / $1) as database/sql.
// Additionally, it allows you to bind other fragment builders.
func F(raw string, args ...any) *Fragment {
	return &Fragment{
		raw:  raw,
		args: args,
	}
}

// WithArgs replaces the args of f.
func (f *Fragment) WithArgs(args ...any) *Fragment {
	f.args = args
	return f
}

// AppendArgs appends args to f.
func (f *Fragment) AppendArgs(args ...any) *Fragment {
	f.args = append(f.args, args...)
	return f
}

// SetArg sets the arg at the given index (0-based, -1 means last).
func (f *Fragment) SetArg(index int, args any) *Fragment {
	if index < 0 {
		index = len(f.args) + index
	}
	if index < 0 || index >= len(f.args) {
		panic(fmt.Errorf("index %d out of range [%d,%d)", index, -len(f.args), len(f.args)))
	}
	f.args[index] = args
	return f
}

// NoUsageCheck disables usage checking whether all bind vars are used.
// This is useful when the arguments are constructed not by hand but by programs dynamically.
func (f *Fragment) NoUsageCheck() *Fragment {
	f.noUsageCheck = true
	return f
}
