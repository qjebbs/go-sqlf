package sqlf

var _ FragmentBuilder = (*arg)(nil)

type arg struct {
	any
}

// BuildFragment implements FragmentBuilder
func (c *arg) BuildFragment(ctx *Context) (query string, err error) {
	built := ctx.CommitArg(c.any)
	return built, nil
}
