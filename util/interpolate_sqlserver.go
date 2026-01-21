package util

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/qjebbs/go-sqlf/v4/dialect"
)

var _ Dialect = _SQLServer{}

type _SQLServer struct {
	dialect.SQLServer
}

// FormatTime formats time strings for the dialect.
func (d _SQLServer) FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}

func (d _SQLServer) QuoteString(s string) string {
	if s == "" {
		return "''"
	}
	var hasUnicode bool
	for _, r := range s {
		if r > 255 {
			hasUnicode = true
			break
		}
	}

	const concatOperator = " + "
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
			if r > 255 {
				b.WriteString("NCHAR(")
			} else {
				b.WriteString("CHAR(")
			}
			b.WriteString(strconv.FormatInt(int64(r), 10))
			b.WriteString(")")
		} else {
			if !inQuote {
				if hasUnicode {
					b.WriteString("N'")
				} else {
					b.WriteString("'")
				}
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
