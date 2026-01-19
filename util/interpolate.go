package util

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/qjebbs/go-sqlf/v4/dialect"
	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

// interpolateDefaultDialect is a dialect used for interpolation when no dialect is provided.
// It uses ANSI SQL as the base dialect, but not specific about bind variable style,
// leaving it to be determined by the parser.
type interpolateDefaultDialect struct {
	dialect.AnsiSQL
}

func (d interpolateDefaultDialect) BindVarStyle() dialect.BindVarStyle {
	return dialect.BindVarStyleDefault // syntax.BindVarStyleUnknown
}

// Interpolate interpolates the args into the query.
//
// The dialect parameter is optional. If provided, it will be used to
// determine the bind variable style and to format time values.
//
// !!! Use it only on debug purposes. Due to complexity of SQL syntax
// between different dialects, it may not work correctly in all cases.
func Interpolate(query string, args []any, d ...dialect.Dialect) (string, bool) {
	var dialect dialect.Dialect
	if len(d) > 0 && d[0] != nil {
		dialect = d[0]
	} else {
		dialect = interpolateDefaultDialect{}
	}
	ok := true
	exprs, err := syntax.ParseForInterpolating(query, dialect.BindVarStyle())
	if err != nil {
		ok = false
	}
	b := new(strings.Builder)
	for _, decl := range exprs.ExprList {
		switch decl := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(decl.Text)
		case *syntax.BindVarExpr:
			if decl.Name != "" {
				val, found := findNamedArg(decl.Name, args)
				if !found {
					ok = false
					b.Write([]byte(fmt.Sprintf("/* %s not found */", decl.Name)))
				} else {
					v, err := encodeValue(dialect, val)
					if err != nil {
						ok = false
						v = []byte(fmt.Sprintf("/* %s */", err.Error()))
					}
					b.Write(v)
				}
				continue
			}
			if decl.Index < 1 || decl.Index > len(args) {
				ok = false
				b.Write([]byte(fmt.Sprintf("/* bindvar index out of range: %d */", decl.Index)))
				continue
			}
			v, err := encodeValue(dialect, args[decl.Index-1])
			if err != nil {
				ok = false
				v = []byte(fmt.Sprintf("/* %s */", err.Error()))
			}
			b.Write(v)
		default:
			ok = false
			b.Write([]byte(fmt.Sprintf("/* unsupported declaration: %T */", decl)))
		}
	}
	return b.String(), ok
}

func findNamedArg(name string, args []any) (any, bool) {
	for _, arg := range args {
		switch v := arg.(type) {
		case sql.NamedArg:
			if v.Name == name {
				return v.Value, true
			}
		}
	}
	return nil, false
}

func encodeValue(dialect dialect.Dialect, arg any) ([]byte, error) {
	if arg == nil {
		return []byte("NULL"), nil
	}
	rv := reflect.ValueOf(arg)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return []byte("NULL"), nil
		}
	}
	buf := bytes.NewBuffer(nil)
	switch v := arg.(type) {
	case driver.Valuer:
		val, err := v.Value()
		if err != nil {
			return nil, err
		}
		enc, err := encodeValue(dialect, val)
		if err != nil {
			return nil, err
		}
		buf.Write(enc)
	case time.Time:
		timeFormat := dialect.TimeFormat()
		if !strings.Contains(timeFormat, "07") {
			v = v.UTC()
		}
		v = v.Round(time.Microsecond)
		buf.WriteRune('\'')
		buf.WriteString(v.Format(timeFormat))
		buf.WriteRune('\'')
	case fmt.Stringer:
		buf.Write(quoteStringValue(v.String()))
	default:
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return []byte("NULL"), nil
			}
			rv = rv.Elem()
		}
		switch k := rv.Kind(); k {
		case reflect.Bool:
			if rv.Bool() {
				buf.WriteString("TRUE")
			} else {
				buf.WriteString("FALSE")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString(fmt.Sprintf("%d", rv.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			buf.WriteString(fmt.Sprintf("%d", rv.Uint()))
		case reflect.Float32, reflect.Float64:
			buf.WriteString(fmt.Sprintf("%f", rv.Float()))
		case reflect.String:
			buf.Write(quoteStringValue(rv.String()))
		default:
			return nil, fmt.Errorf("unsupported type %T", arg)
		}
	}
	return buf.Bytes(), nil
}

func quoteStringValue(s string) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('\'')
	buf.WriteString(strings.ReplaceAll(s, "'", "''"))
	buf.WriteRune('\'')
	return buf.Bytes()
}

var escaping = []struct {
	from rune
	to   string
}{
	{'\x00', `\0`},
	{'\n', `\n`},
	{'\r', `\r`},
	{'\b', `\b`},
	{'\t', `\t`},
	{'\x1a', `\Z`},
	{'\'', "''"},
	{'"', `\"`},
	{'\\', `\\`},
}
