package sqlf

// Properties is the Properties
type Properties struct {
	Args      *ArgsProperty
	Columns   *ColumnsProperty
	Tables    *TablesProperty
	Fragments *FragmentsProperty
	Builders  *BuildersProperty
}

// NewProperties returns a new Properties.
func NewProperties(f *Fragment) *Properties {
	if f == nil {
		return nil
	}
	return &Properties{
		Args:      NewArgsProperty(f.Args),
		Columns:   NewColumnsProperty(f.Columns),
		Tables:    NewTablesProperty(f.Tables),
		Fragments: NewFragmentsProperty(f.Fragments),
		Builders:  NewBuildersProperty(f.Builders),
	}
}
