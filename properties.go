package sqlf

// properties is the properties
type properties struct {
	Args      *ArgsProperty
	Columns   *ColumnsProperty
	Tables    *TablesProperty
	Fragments *FragmentsProperty
	Builders  *BuildersProperty
}

// newProperties returns a new Properties.
func newProperties(f *Fragment) *properties {
	if f == nil {
		return nil
	}
	return &properties{
		Args:      NewArgsProperty(f.Args),
		Columns:   NewColumnsProperty(f.Columns),
		Tables:    NewTablesProperty(f.Tables),
		Fragments: NewFragmentsProperty(f.Fragments),
		Builders:  NewBuildersProperty(f.Builders),
	}
}
