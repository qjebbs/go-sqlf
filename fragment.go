package sqlf

// Fragment is the builder for a part of or even the full query, it allows you
// to write and combine fragments with freedom.
type Fragment struct {
	Raw       string            // Raw string support bind vars (?, $1) and preprocessing functions (#join).
	Args      []any             // Args can be referenced by the Raw, for example: ?, $1
	Fragments []FragmentBuilder // Fragments can be referenced by the Raw, for example: #f1, #fragment1
	Prefix    string            // Prefix is added before the fragment only when the fragment is built not empty.
	Suffix    string            // Suffix is added after the fragment only when the fragment is built not empty.
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

// WithFragments sets the fragments of f.
func (f *Fragment) WithFragments(fragments ...FragmentBuilder) *Fragment {
	f.Fragments = fragments
	return f
}

// AppendArgs appends args to f.
func (f *Fragment) AppendArgs(args ...any) *Fragment {
	f.Args = append(f.Args, args...)
	return f
}

// AppendFragments appends fragments to f.
func (f *Fragment) AppendFragments(fragments ...FragmentBuilder) *Fragment {
	f.Fragments = append(f.Fragments, fragments...)
	return f
}
