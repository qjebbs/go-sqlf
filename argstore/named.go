package argstore

import (
	"database/sql"
	"reflect"
	"strconv"
)

var _ Store = (*Named)(nil)

// Named implements ArgStore for Named parameters (e.g., ":name", "@name").
type Named struct {
	marker string
	prefix string
	args   []any
	dict   map[any]int
}

// NewNamed creates a new Named ArgStore.
func NewNamed(marker, prefix string) *Named {
	return &Named{
		marker: marker,
		prefix: prefix,
		dict:   make(map[any]int),
	}
}

// Args implements ArgStore.Args.
func (s *Named) Args() []any {
	return s.args
}

// CommitArg implements ArgStore.CommitArg.
func (s *Named) CommitArg(arg any) string {
	if arg != nil && reflect.TypeOf(arg).Comparable() {
		if i, ok := s.dict[arg]; ok {
			return s.marker + strconv.Itoa(i)
		}
	}
	i := len(s.args) + 1
	s.dict[arg] = i
	baseName := s.prefix + strconv.Itoa(i)
	name := s.marker + baseName
	s.args = append(s.args, sql.Named(baseName, arg))
	return name
}
