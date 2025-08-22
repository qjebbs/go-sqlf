package util

import (
	"github.com/qjebbs/go-sqlf/v2"
	"github.com/qjebbs/go-sqlf/v2/syntax"
)

// Build builds the SQL query and arguments from the given fragment builder.
func Build(b sqlf.Builder, bindVarStyle syntax.BindVarStyle) (query string, args []any, err error) {
	ctx := sqlf.NewContext(bindVarStyle)
	query, err = b.BuildFragment(ctx)
	if err != nil {
		return "", nil, err
	}
	args = ctx.Args()
	return query, args, nil
}
