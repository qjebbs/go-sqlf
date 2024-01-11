package sqls

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
func (s *Fragment) AppendTables(tables ...Table) {
	s.Tables = append(s.Tables, tables...)
}

// AppendColumns appends columns to the fragment.
func (s *Fragment) AppendColumns(columns ...*TableColumn) {
	s.Columns = append(s.Columns, columns...)
}

// AppendFragments appends fragments to the s.Fragments.
func (s *Fragment) AppendFragments(fragments ...*Fragment) {
	s.Fragments = append(s.Fragments, fragments...)
}

// AppendArgs appends args to the s.Args.
func (s *Fragment) AppendArgs(args ...any) {
	s.Args = append(s.Args, args...)
}

// WithTables replace s.Tables with the tables
func (s *Fragment) WithTables(tables ...Table) {
	s.Tables = tables
}

// WithColumns replace s.Columns with the columns
func (s *Fragment) WithColumns(columns ...*TableColumn) {
	s.Columns = columns
}

// WithFragments replace s.Fragments with the fragments
func (s *Fragment) WithFragments(fragments ...*Fragment) {
	s.Fragments = fragments
}

// WithArgs replace s.Args with the args
func (s *Fragment) WithArgs(args ...any) {
	s.Args = args
}
