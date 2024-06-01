package sqlf

import (
	"strconv"
)

// Args returns the built args of the context.
func (c *Context) Args() []any {
	return c.root().argStore.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
//
// It's used usually in the implementation of a FragmentBuilder,
// most users don't need to care about it.
func (c *Context) CommitArg(arg any) string {
	return c.root().argStore.CommitArg(arg)
}

type argStore interface {
	Args() []any
	CommitArg(arg any) string
}

type questionArgStore struct {
	args []any
}

func newQuestionArgStore() *questionArgStore {
	return &questionArgStore{}
}

func (s *questionArgStore) Args() []any {
	return s.args
}

func (s *questionArgStore) CommitArg(arg any) string {
	s.args = append(s.args, arg)
	return "?"
}

type dollarArgStore struct {
	args []any
	dict map[any]int
}

func newDollarArgStore() *dollarArgStore {
	return &dollarArgStore{
		dict: make(map[any]int),
	}
}

func (s *dollarArgStore) Args() []any {
	return s.args
}

func (s *dollarArgStore) CommitArg(arg any) string {
	if i, ok := s.dict[arg]; ok {
		return "$" + strconv.Itoa(i)
	}
	i := len(s.args) + 1
	s.dict[arg] = i
	s.args = append(s.args, arg)
	return "$" + strconv.Itoa(i)
}
