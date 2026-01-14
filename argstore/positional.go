package argstore

var _ Store = (*Positional)(nil)

// Positional implements Store for positional parameters (e.g., "?").
type Positional struct {
	args []any
}

// NewPositional creates a new Positional Store.
func NewPositional() *Positional {
	return &Positional{}
}

// Args implements Store.Args.
func (p *Positional) Args() []any {
	return p.args
}

// CommitArg implements Store.CommitArg.
func (p *Positional) CommitArg(arg any) string {
	p.args = append(p.args, arg)
	return "?"
}
