package sqlf

// Fa creates a new Fragment with Args property.
func Fa(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// Ff creates a new Fragment with Fragments property.
func Ff(raw string, fragments ...*Fragment) *Fragment {
	return &Fragment{
		Raw:       raw,
		Fragments: fragments,
	}
}

// Fb creates a new Fragment with Builders property.
func Fb(raw string, builders ...Builder) *Fragment {
	return &Fragment{
		Raw:      raw,
		Builders: builders,
	}
}

// Fc creates a new Fragment with Columns property.
func Fc(raw string, columns ...*Column) *Fragment {
	return &Fragment{
		Raw:     raw,
		Columns: columns,
	}
}

// Ft creates a new Fragment with Tables property.
func Ft(raw string, tables ...Table) *Fragment {
	return &Fragment{
		Raw:    raw,
		Tables: tables,
	}
}
