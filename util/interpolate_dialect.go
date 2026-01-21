package util

import (
	"time"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

// Dialect defines an SQL dialect that supports string interpolation.
// It extends the base dialect.Dialect with string quoting capabilities.
type Dialect interface {
	dialect.Dialect
	// FormatTime formats time strings for the dialect.
	FormatTime(t time.Time) string
	// QuoteString quotes a string for use in a query.
	// It's used in Interpolate() only.
	//
	// It should escape any special characters as needed by the dialect,
	// and take care of any necessary prefixing (e.g. E'...' in PostgreSQL, N'...' in SQL Server).
	//
	// Examples:
	//   QuoteString("str") // 'str'
	//   SQLite.QuoteString("str\nstr") // 'str' || CHAR(10) || 'str'
	//   PostgreSQL.QuoteString("str\nstr") // E'str\nstr'
	QuoteString(s string) string
}

// defaultDialect is a dialect used for interpolation when no dialect is provided.
// It uses ANSI SQL as the base dialect, but not specific about bind variable style,
// leaving it to be determined by the parser.
type defaultDialect struct {
	_AnsiSQL
}

func (d defaultDialect) BindVarStyle() dialect.BindVarStyle {
	return dialect.BindVarStyleDefault // syntax.BindVarStyleUnknown
}

func asMyDialect(d dialect.Dialect) Dialect {
	if iq, ok := d.(Dialect); ok {
		return iq
	}
	switch v := d.(type) {
	case dialect.AnsiSQL:
		return _AnsiSQL{v}
	case dialect.SQLite:
		return _SQLite{v}
	case dialect.PostgreSQL:
		return _PostgreSQL{v}
	case dialect.MySQL:
		return _MySQL{v}
	case dialect.SQLServer:
		return _SQLServer{v}
	case dialect.Oracle:
		return _Oracle{v}
	default:
		// Fallback to ANSI SQL dialect.
		return _Other{Dialect: d}
	}
}

var _ Dialect = _Other{}

type _Other struct {
	dialect.Dialect
	ansi _AnsiSQL
}

func (d _Other) QuoteString(s string) string {
	return d.ansi.QuoteString(s)
}

func (d _Other) FormatTime(t time.Time) string {
	return d.ansi.FormatTime(t)
}
