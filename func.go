package sqlf

import (
	"fmt"
	"reflect"
	"unicode"
)

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	stringType       = reflect.TypeOf((*string)(nil)).Elem()
	fmtStringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()

	contextPointerType = reflect.TypeOf((*context)(nil))
)

// FuncMap is the type of the map defining the mapping from names to functions.
type FuncMap map[string]any

// createValueFuncs turns a FuncMap into a map[string]reflect.Value
func createValueFuncs(funcMap FuncMap) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	addValueFuncs(m, funcMap)
	return m
}

// addValueFuncs adds to values the functions in funcs, converting them to reflect.Values.
func addValueFuncs(out map[string]reflect.Value, in FuncMap) {
	for name, fn := range in {
		if !goodName(name) {
			panic(fmt.Errorf("function name %q is not a valid identifier", name))
		}
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			panic("value for " + name + " not a function")
		}
		if !goodFunc(v.Type()) {
			panic(fmt.Errorf("can't install method/function %q with %d results", name, v.Type().NumOut()))
		}
		out[name] = v
	}
}

// goodName reports whether the function name is a valid identifier.
func goodName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		switch {
		case r == '_':
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			return false
		}
	}
	return true
}

// goodFunc reports whether the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
	// allow functions:
	//  func(...) string
	//  func(...) (string, error)
	switch typ.NumOut() {
	case 1:
		return typ.Out(0) == stringType
	case 2:
		return typ.Out(0) == stringType && typ.Out(1) == errorType
	}
	return false
}
