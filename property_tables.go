package sqlf

// TablesProperty is the tables property
type TablesProperty struct {
	*propertyBase[Table]
}

// NewTablesProperty returns a new TablesProperty.
func NewTablesProperty(tables []Table) *TablesProperty {
	return &TablesProperty{
		propertyBase: newPropertyBase("tables", tables),
	}
}

// Build builds the arg at index, with cache.
func (b *TablesProperty) Build(ctx *Context, index int) (string, error) {
	if err := b.validateIndex(index); err != nil {
		return "", err
	}
	b.used[index-1] = true
	return string(b.items[index-1]), nil
}
