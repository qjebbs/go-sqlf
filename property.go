package sqlf

import (
	"fmt"
	"strings"
)

type property[T any] struct {
	name  string
	items []T
	cache []string
	used  []bool
}

// newProperty returns a new property.
func newProperty[T any](name string, items []T) *property[T] {
	return &property[T]{
		name:  name,
		items: items,
		cache: make([]string, len(items)),
		used:  make([]bool, len(items)),
	}
}

// Count returns the count of items.
func (b *property[T]) Count() int {
	return len(b.items)
}

// Get returns the item at index, starting from 1.
func (b *property[T]) Get(index int) (T, error) {
	var zero T
	if err := b.validateIndex(index); err != nil {
		return zero, err
	}
	return b.items[index-1], nil
}

// ReportUsed reports the item at index is used, starting from 1.
func (b *property[T]) ReportUsed(index int) {
	if index < 1 || index > len(b.items) {
		return
	}
	b.used[index-1] = true
}

// CheckUsage checks if all items are used.
func (b *property[T]) CheckUsage() error {
	unused := new(strings.Builder)
	for i, v := range b.used {
		if !v {
			if unused.Len() > 0 {
				unused.WriteString(", ")
			}
			unused.WriteString(fmt.Sprint(i + 1))
		}
	}
	if unused.Len() == 0 {
		return nil
	}
	return fmt.Errorf("%s not used: [%s]", b.name, unused.String())
}

// validateIndex validates the index.
func (b *property[T]) validateIndex(index int) error {
	if index < 1 || index > len(b.items) {
		return fmt.Errorf("%w: %s index %d out of range [1,%d]", ErrInvalidIndex, b.name, index, len(b.items))
	}
	return nil
}
