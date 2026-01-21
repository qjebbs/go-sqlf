package util

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ Dialect = _Oracle{}

type _Oracle struct {
	dialect.Oracle
}

// FormatTime formats time strings for the dialect.
func (d _Oracle) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

func (d _Oracle) QuoteString(s string) string {
	if s == "" {
		return "''"
	}
	const concatOperator = " || "
	var b strings.Builder
	inQuote := false
	for i, r := range s {
		if i > 0 && !inQuote {
			b.WriteString(concatOperator)
		}
		isControl := unicode.IsControl(r)
		if isControl {
			if inQuote {
				b.WriteString("'")
				b.WriteString(concatOperator)
				inQuote = false
			}
			b.WriteString("CHR(")
			b.WriteString(strconv.FormatInt(int64(r), 10))
			b.WriteString(")")
		} else {
			if !inQuote {
				b.WriteString("'")
				inQuote = true
			}
			if r == '\'' {
				b.WriteString("''")
			} else {
				b.WriteRune(r)
			}
		}
	}
	if inQuote {
		b.WriteString("'")
	}
	return b.String()
}
