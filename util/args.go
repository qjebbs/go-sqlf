package util

import (
	"reflect"
)

// ArgsFlatted is a help func to create a slice of query args.
//
// it accepts values, slices and arrays, and the all 1st depth elements
// in the slice/array will be extract and concatenated to the returned
// slice, which is called flattening, e.g.:
//
//	ArgsFlatted(1, []int{2, 3}, []string{"a", "b", "c"}) => []any{1, 2, 3, "a", "b", "c"}
//	ArgsFlatted(1, []int{}, []int{2}) => []any{1, 2}
func ArgsFlatted(valueOrSlices ...any) []any {
	args := make([]any, 0, 10)
	for _, v := range valueOrSlices {
		args = append(args, argsFrom(v)...)
	}
	return args
}

func argsFrom(v any) []any {
	switch a := v.(type) {
	case []any:
		return a
	case []bool:
		return Ttoa(a)
	case []float64:
		return Ttoa(a)
	case []float32:
		return Ttoa(a)
	case []int64:
		return Ttoa(a)
	case []int32:
		return Ttoa(a)
	case []int:
		return Ttoa(a)
	case []uint64:
		return Ttoa(a)
	case []uint32:
		return Ttoa(a)
	case []uint:
		return Ttoa(a)
	case []string:
		return Ttoa(a)
	case *[]bool:
		return Ttoa(*a)
	case *[]float64:
		return Ttoa(*a)
	case *[]float32:
		return Ttoa(*a)
	case *[]int64:
		return Ttoa(*a)
	case *[]int32:
		return Ttoa(*a)
	case *[]int:
		return Ttoa(*a)
	case *[]uint64:
		return Ttoa(*a)
	case *[]uint32:
		return Ttoa(*a)
	case *[]uint:
		return Ttoa(*a)
	case *[]string:
		return Ttoa(*a)
	default:
		return convertArrayReflect(v)
	}
}

// Ttoa is a help func to convert a slice to []any.
func Ttoa[T any](slice []T) []any {
	if len(slice) == 0 {
		return nil
	}
	b := make([]any, 0, len(slice))
	for _, v := range slice {
		b = append(b, v)
	}
	return b
}

func convertArrayReflect(slice any) []any {
	rv := reflect.ValueOf(slice)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return nil
		}
	case reflect.Array:
	default:
		return []any{slice}
	}
	n := rv.Len()
	if n == 0 {
		return nil
	}
	s := make([]any, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, rv.Index(i).Interface())
	}
	return s
}
