package util

import (
	"strings"
	"time"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ Dialect = _AnsiSQL{}

type _AnsiSQL struct {
	dialect.AnsiSQL
}

// FormatTime formats time strings for the dialect.
func (d _AnsiSQL) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

func (d _AnsiSQL) QuoteString(s string) string {
	var b strings.Builder
	b.WriteString("'")
	for _, r := range s {
		if r == '\x00' {
			// Prevent query truncation attack,
			// repsect it as end of string.
			break
		}
		if r == '\'' {
			b.WriteString("''")
			continue
		}
		b.WriteRune(r)
	}
	b.WriteString("'")
	return b.String()
}
