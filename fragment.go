package sqlf

// Fragment is the builder for a part of or even the full query, it allows you
// to write and combine fragments with freedom.
type Fragment struct {
	Raw       string      // Raw string support bind vars and preprocessing functions.
	Args      []any       // Args to be referred by the Raw, e.g.: ?, $1
	Columns   []*Column   // Columns to be referred by the Raw, e.g.: #c1, #column2
	Tables    []Table     // Table names / alias to be referred by the Raw, e.g.: #t1, #table2
	Fragments []*Fragment // Fragments to be referred by the Raw, e.g.: #fragment1, #fragment2
	Builders  []Builder   // Builders to be referred by the Raw, e.g.: #builder1, #builder2

	Prefix string // Prefix is added before the built fragment only if which is not empty.
	Suffix string // Suffix is added after the built fragment only if which is not empty.
}

// AppendArgs appends args to the fragment.
// Args are used to be referred by the Raw, e.g.: ?, $1
func (f *Fragment) AppendArgs(args ...any) {
	f.Args = append(f.Args, args...)
}

// AppendColumns appends columns to the fragment.
// Columns are used to be referred by the Raw, e.g.: #c1, #column2
func (f *Fragment) AppendColumns(columns ...*Column) {
	f.Columns = append(f.Columns, columns...)
}

// AppendTables appends tables to the fragment.
// Tables are used to be referred by the Raw, e.g.: #t1, #table2
func (f *Fragment) AppendTables(tables ...Table) {
	f.Tables = append(f.Tables, tables...)
}

// AppendFragments appends fragments to the fragment.
// Fragments are used to be referred by the Raw, e.g.: #fragment1, #fragment2
func (f *Fragment) AppendFragments(fragments ...*Fragment) {
	f.Fragments = append(f.Fragments, fragments...)
}

// AppendBuilders appends builders to the fragment.
// Builders are used to be referred by the Raw, e.g.: #builder1, #builder2
func (f *Fragment) AppendBuilders(builders ...Builder) {
	f.Builders = append(f.Builders, builders...)
}

// WithArgs replace f.Args with the args
// Args are used to be referred by the Raw, e.g.: ?, $1
func (f *Fragment) WithArgs(args ...any) {
	f.Args = args
}

// WithColumns replace f.Columns with the columns
// Columns are used to be referred by the Raw, e.g.: #c1, #column2
func (f *Fragment) WithColumns(columns ...*Column) {
	f.Columns = columns
}

// WithTables replace f.Tables with the tables
// Tables are used to be referred by the Raw, e.g.: #t1, #table2
func (f *Fragment) WithTables(tables ...Table) {
	f.Tables = tables
}

// WithFragments replace f.Fragments with the fragments
// Fragments are used to be referred by the Raw, e.g.: #fragment1, #fragment2
func (f *Fragment) WithFragments(fragments ...*Fragment) {
	f.Fragments = fragments
}

// WithBuilders replace f.Builders with the builders
// Builders are used to be referred by the Raw, e.g.: #builder1, #builder2
func (f *Fragment) WithBuilders(builders ...Builder) {
	f.Builders = builders
}
