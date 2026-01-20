package sqlf

import "fmt"

// Fragment represents a composable SQL query fragment.
//
// Use the F() function to create a new Fragment.
type Fragment struct {
	raw          string // Raw string support bind vars (?, $1)
	args         []any  // Args that can be referenced by the Raw. An arg can be either an ordinary arg or a Builder.
	noUsageCheck bool   // If true, skip checking whether all bind vars are used.
}

// F creates a new Fragment.
//
// # Bind Variables
//
// Placeholders in the raw SQL string can be used to bind arguments.
// These arguments can be simple values or other fragment builders
// who implement sqlf.Builder.
//
// There are two types of placeholders:
//   - `?`: A positional placeholder. Arguments are bound in order.
//   - `$N`: An indexed placeholder (e.g., $1, $2), where N is a 1-based index.
//
// To include a literal `?` or `$` in your SQL, double it (e.g., `??` or `$$`).
//
//	// Bind a simple value using a positional placeholder
//	cond := sqlf.F("name = ?", "jebbs")
//	// Bind another builder
//	where := sqlf.F("WHERE ?", cond)
//
//	// Escape a placeholder to include it literally in the output.
//	// This example builds "WHERE foo -> $1 ? $2" for PostgreSQL.
//	sqlf.F("WHERE foo -> $1 ?? $2", "bar", "baz")
//
// # Identifiers and Literals
//
// SQL identifiers and string literals can be quoted using single quotes ('string'),
// double quotes ("identifier"), or backticks (`identifier`).
// Placeholders within quoted sections are ignored.
//
// To include a quote character within a quoted string, double it.
//
//	// Builds `[name] = 'alice''s friend'` for SQL Server
//	sqlf.F(`"name" = 'alice''s friend'`)
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
