package util

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ Dialect = _PostgreSQL{}

type _PostgreSQL struct {
	dialect.PostgreSQL
}

// FormatTime formats time strings for the dialect.
func (d _PostgreSQL) FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999-07:00")
}

func (d _PostgreSQL) QuoteString(s string) string {
	var b strings.Builder
	var hasEscape bool
	b.WriteString("'")
	for _, r := range s {
		if r == '\'' {
			b.WriteString("''")
			continue
		}
		isControl := unicode.IsControl(r)
		if isControl {
			hasEscape = true
		}
		if !isControl {
			b.WriteRune(r)
		} else {
			switch r {
			case '\b':
				b.WriteString(`\b`)
			case '\f':
				b.WriteString(`\f`)
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			case '\t':
				b.WriteString(`\t`)
			default:
				if r < 256 {
					b.WriteString(fmt.Sprintf("\\x%02x", r))
				} else {
					b.WriteString(fmt.Sprintf("\\u%04x", r))
				}
			}
		}
	}
	b.WriteString("'")
	if hasEscape {
		return "E" + b.String()
	}
	// No escape sequences, return simple quoted string
	return b.String()
}
