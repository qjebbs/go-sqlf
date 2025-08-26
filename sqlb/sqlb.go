// Package sqlb provides a complex SQL query builder shipped with
// WITH-CTE / JOIN Elimination capabilities, while *sqlf.Fragment
// is the underlying foundation.
package sqlb

import (
	"github.com/qjebbs/go-sqlf/v3/syntax"
)

// Builder is the interface for sql builders.
type Builder interface {
	// BuildQuery builds and returns the query and args.
	BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error)
}
