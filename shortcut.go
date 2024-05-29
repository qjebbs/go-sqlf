package sqlf

// F creates a new Fragment.
func F(raw string) *Fragment {
	return &Fragment{
		Raw: raw,
	}
}

// Fa creates a new Fragment with args.
func Fa(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// Ff creates a new Fragment with fragments.
func Ff(raw string, fragments ...FragmentBuilder) *Fragment {
	return &Fragment{
		Raw:       raw,
		Fragments: fragments,
	}
}
