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
//
// # Bind Variables
//
// ? / $1 in the raw string can refer to both ordinary args and fragment builders in args.
// To use them as ordinary characters outside quotes, double them.
//
//	// refer to an ordinary arg
//	cond := sqlf.F("name = ?", "jebbs")
//	// refer to a fragment builder
//	sqlf.F("WHERE ?", cond)
//	// will be built as "WHERE foo -> $1 ? $2" for PostgreSQL
//	sqlf.F("WHERE foo -> $1 ?? $2", "bar", "baz")
//
// # Identifiers and Literals
//
// Single quotes ('string'), double quotes ("name"), and backticks (`name`)
// are used for quoting. Bind variables inside quoted strings are not
// processed. To include a quote character within a quoted string, double it.
//
//	// will be built as "[name] = 'alice'" for SQL Server
//	sqlf.F(`"name" = 'alice'`)
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
