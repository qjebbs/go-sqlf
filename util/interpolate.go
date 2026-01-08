package util

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/qjebbs/go-sqlf/v4/internal/syntax"
)

// InterpolateOption is the option of Interpolate.
type InterpolateOption func(*interpolateOptions)

type interpolateOptions struct {
	// TimeFormat is the format of time value.
	TimeFormat string
}

func defaultInterpolateOptions() *interpolateOptions {
	return &interpolateOptions{
		TimeFormat: time.RFC3339,
	}
}

func applyInterpolateOptions(options []InterpolateOption) *interpolateOptions {
	opts := defaultInterpolateOptions()
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

// WithInterpolateTimeFormat sets the format of time value.
func WithInterpolateTimeFormat(format string) InterpolateOption {
	return func(opts *interpolateOptions) {
		opts.TimeFormat = format
	}
}

// Interpolate interpolates the args into the query.
//
// !!! Use it only on debug purposes.
func Interpolate(query string, args []any, options ...InterpolateOption) (string, error) {
	opts := applyInterpolateOptions(options)
	exprs, err := syntax.Parse(query)
	if err != nil {
		return "", err
	}
	b := new(strings.Builder)
	for _, decl := range exprs.ExprList {
		switch decl := decl.(type) {
		case *syntax.PlainExpr:
			b.WriteString(decl.Text)
		case *syntax.BindVarExpr:
			v, err := encodeValue(args[decl.Index-1], opts)
			if err != nil {
				return "", err
			}
			b.Write(v)
		default:
			return "", fmt.Errorf("%s: unsupported declaration", decl.Pos())
		}
	}
	return b.String(), nil
}

func encodeValue(arg any, opts *interpolateOptions) ([]byte, error) {
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
		enc, err := encodeValue(val, opts)
		if err != nil {
			return nil, err
		}
		buf.Write(enc)
	case time.Time:
		if v.IsZero() {
			buf.WriteString("'0000-00-00'")
			break
		}
		// In SQL standard, the precision of fractional seconds in time literal is up to 6 digits.
		v = v.Round(time.Microsecond)
		buf.WriteRune('\'')
		buf.WriteString(v.Format(opts.TimeFormat))
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
