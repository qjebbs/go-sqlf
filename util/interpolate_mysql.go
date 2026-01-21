package util

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ Dialect = _MySQL{}

type _MySQL struct {
	dialect.MySQL
}

// FormatTime formats time strings for the dialect.
func (d _MySQL) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

func (d _MySQL) QuoteString(s string) string {
	var b strings.Builder
	b.WriteString("'")
	for _, r := range s {
		if r == '\'' {
			b.WriteString("''")
			continue
		}
		isControl := unicode.IsControl(r)
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
	return b.String()
}
