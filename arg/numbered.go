package arg

import (
	"reflect"
	"strconv"
)

var _ Store = (*Numbered)(nil)

// Numbered implements Store for Numbered parameters (e.g., "?1", "?2").
type Numbered struct {
	marker string
	args   []any
	dict   map[any]int
}

// NewNumbered creates a new Numbered Store.
func NewNumbered(marker string) *Numbered {
	return &Numbered{
		marker: marker,
		dict:   make(map[any]int),
	}
}

// Args implements Store.Args.
func (s *Numbered) Args() []any {
	return s.args
}

// CommitArg implements Store.CommitArg.
func (s *Numbered) CommitArg(arg any) string {
	if arg != nil && reflect.TypeOf(arg).Comparable() {
		if i, ok := s.dict[arg]; ok {
			return s.marker + strconv.Itoa(i)
		}
	}
	i := len(s.args) + 1
	s.dict[arg] = i
	s.args = append(s.args, arg)
	return s.marker + strconv.Itoa(i)
}
