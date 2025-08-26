// Package sqlb is a SQL query builder based on `sqlf.Fragment`.
package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

// Builder is the interface for sql builders.
type Builder interface {
	sqlf.Builder

	// BuildQuery builds and returns the query and args.
	BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error)
}
