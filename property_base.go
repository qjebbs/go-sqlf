package sqlf

import (
	"fmt"
	"strings"
)

type propertyBase[T any] struct {
	name  string
	items []T
	cache []string
	used  []bool
}

// newPropertyBase returns a new ArgsBuilder.
func newPropertyBase[T any](name string, items []T) *propertyBase[T] {
	return &propertyBase[T]{
		name:  name,
		items: items,
		cache: make([]string, len(items)),
		used:  make([]bool, len(items)),
	}
}

// ReportUsed reports the arg at index is used.
func (b *propertyBase[T]) ReportUsed(index int) {
	if index < 1 || index > len(b.items) {
		return
	}
	b.used[index-1] = true
}

// CheckUsage checks if all args are used.
func (b *propertyBase[T]) CheckUsage() error {
	if msg := unusedIndexs(b.used); msg != "" {
		return fmt.Errorf("%s not used: [%s]", b.name, msg)
	}
	return nil
}

// validateIndex validates the index.
func (b *propertyBase[T]) validateIndex(index int) error {
	if index < 1 || index > len(b.items) {
		return fmt.Errorf("%w: %s index %d out of range [1,%d]", ErrInvalidIndex, b.name, index, len(b.items))
	}
	return nil
}

func unusedIndexs(used []bool) string {
	unused := new(strings.Builder)
	for i, v := range used {
		if !v {
			if unused.Len() > 0 {
				unused.WriteString(", ")
			}
			unused.WriteString(fmt.Sprint(i + 1))
		}
	}
	return unused.String()
}
