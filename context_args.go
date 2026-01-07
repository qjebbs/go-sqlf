package sqlf

import (
	"strconv"
)

type argStoreKey struct{}

// Args returns the built args of the context.
func (c *Context) Args() []any {
	store := c.Value(argStoreKey{}).(argStore)
	return store.Args()
}

// CommitArg commits an built arg to the context and returns the built bindvar.
func (c *Context) CommitArg(arg any) string {
	store := c.Value(argStoreKey{}).(argStore)
	return store.CommitArg(arg)
}

func newArgStore(style BindStyle) argStore {
	if style == BindStyleDollar {
		return newDollarArgStore()
	}
	return newQuestionArgStore()
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
