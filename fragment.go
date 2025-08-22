package sqlf

// Fragment is the builder for a part of or even a full query, it allows you
// to write and combine fragments with freedom.
type Fragment struct {
	Raw    string // Raw string support bind vars (?, $1)
	Args   []any  // Args can be referenced by the Raw, for example: ?, $1. An arg can be either a sql arg or fragment builder.
	Prefix string // Prefix is added before the fragment only when the fragment is built not empty.
	Suffix string // Suffix is added after the fragment only when the fragment is built not empty.
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

// WithPrefix sets the prefix which is added before the fragment only when the f is built not empty.
func (f *Fragment) WithPrefix(prefix string) *Fragment {
	f.Prefix = prefix
	return f
}

// WithSuffix sets the suffix which is added before the fragment only when the f is built not empty.
func (f *Fragment) WithSuffix(suffix string) *Fragment {
	f.Suffix = suffix
	return f
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
