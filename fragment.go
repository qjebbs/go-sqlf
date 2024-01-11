package sqlf

// Fragment is the builder for a part of or even the full query, it allows you
// to write and combine fragments with freedom.
type Fragment struct {
	Raw       string         // Raw string support bindvars and preprocessing functions.
	Args      []any          // Args to be referenced by the Raw
	Columns   []*TableColumn // Columns to be referenced by the Raw
	Tables    []Table        // Table names / alias to be referenced by the Raw
	Fragments []*Fragment    // Fragments to be referenced by the Raw
	Builders  []Builder      // Builders to be referenced by the Raw

	Prefix string // Prefix is added before the rendered fragment only if which is not empty.
	Suffix string // Suffix is added after the rendered fragment only if which is not empty.
}

// AppendTables appends tables to the fragment.
func (f *Fragment) AppendTables(tables ...Table) {
	f.Tables = append(f.Tables, tables...)
}

// AppendColumns appends columns to the fragment.
func (f *Fragment) AppendColumns(columns ...*TableColumn) {
	f.Columns = append(f.Columns, columns...)
}

// AppendFragments appends fragments to the fragment.
func (f *Fragment) AppendFragments(fragments ...*Fragment) {
	f.Fragments = append(f.Fragments, fragments...)
}

// AppendArgs appends args to the fragment.
func (f *Fragment) AppendArgs(args ...any) {
	f.Args = append(f.Args, args...)
}

// WithTables replace f.Tables with the tables
func (f *Fragment) WithTables(tables ...Table) {
	f.Tables = tables
}

// WithColumns replace f.Columns with the columns
func (f *Fragment) WithColumns(columns ...*TableColumn) {
	f.Columns = columns
}

// WithFragments replace f.Fragments with the fragments
func (f *Fragment) WithFragments(fragments ...*Fragment) {
	f.Fragments = fragments
}

// WithArgs replace f.Args with the args
func (f *Fragment) WithArgs(args ...any) {
	f.Args = args
}
