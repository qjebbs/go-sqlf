package argstore

var _ Store = (*Positional)(nil)

// Positional implements ArgStore for positional parameters (e.g., "?").
type Positional struct {
	args []any
}

// NewPositional creates a new Positional ArgStore.
func NewPositional() *Positional {
	return &Positional{}
}

// Args implements ArgStore.Args.
func (p *Positional) Args() []any {
	return p.args
}

// CommitArg implements ArgStore.CommitArg.
func (p *Positional) CommitArg(arg any) string {
	p.args = append(p.args, arg)
	return "?"
}
