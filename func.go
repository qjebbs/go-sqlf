package sqlf

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"unicode"
)

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	stringType       = reflect.TypeOf((*string)(nil)).Elem()
	fmtStringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()

	contextPointerType = reflect.TypeOf((*FragmentContext)(nil))
)

// FuncMap is the type of the map defining the mapping from names to functions.
type FuncMap map[string]any

type funcInfo struct {
	name      string
	fn        reflect.Value
	nIn       int            // number of arguments, including the variadic one
	inTypes   []reflect.Type // types of all arguments
	contexArg bool           // if the first argument is *FragmentContext
	nInFixed  int            // number of fixed arguments, except the variadic one
	variadic  bool           // if the last argument is variadic
	nOut      int            // number of outputs
	outTypes  []reflect.Type // types of all outputs

	joinError error // error to return when the function is not compatible with #join()
}

// createValueFuncs turns a FuncMap into a map[string]reflect.Value
func createValueFuncs(funcMap FuncMap) map[string]*funcInfo {
	m := make(map[string]*funcInfo)
	addValueFuncs(m, funcMap)
	return m
}

// addValueFuncs adds to values the functions in funcs, converting them to *funcInfos.
func addValueFuncs(out map[string]*funcInfo, in FuncMap) {
	for name, fn := range in {
		if !goodName(name) {
			panic(fmt.Errorf("function name %q is not a valid identifier", name))
		}
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			panic("value for " + name + " not a function")
		}
		typ := v.Type()

		if _, ok := out[name]; ok {
			panic(fmt.Errorf("function %q already exists", name))
		}
		nIn := typ.NumIn()
		nOut := typ.NumOut()
		fun := &funcInfo{
			name:      name,
			fn:        v,
			nIn:       nIn,
			nInFixed:  nIn,
			nOut:      nOut,
			contexArg: nIn > 0 && typ.In(0) == contextPointerType,
			variadic:  typ.IsVariadic(),
		}
		if fun.variadic {
			fun.nInFixed--
		}
		fun.inTypes = make([]reflect.Type, nIn)
		for i := 0; i < nIn; i++ {
			fun.inTypes[i] = typ.In(i)
		}

		fun.outTypes = make([]reflect.Type, nOut)
		for i := 0; i < nOut; i++ {
			fun.outTypes[i] = typ.Out(i)
		}

		if err := goodFunc(fun); err != nil {
			panic(err)
		}

		if err := joinCompatible(fun); err != nil {
			fun.joinError = err
		}

		out[name] = fun
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
func goodFunc(f *funcInfo) error {
	errInvalidFuncOutput := errors.New("invalid function output, allowed: func(...) string; func(...) (string, error);")
	// Check the result signature.
	switch f.nOut {
	case 1:
		if f.outTypes[0] != stringType {
			return errInvalidFuncOutput
		}
	case 2:
		if f.outTypes[0] != stringType || f.outTypes[1] != errorType {
			return errInvalidFuncOutput
		}
	default:
		return errInvalidFuncOutput
	}

	// Check the argument signature.
	for i, t := range f.inTypes {
		if i == 0 && f.contexArg {
			continue
		}
		if f.variadic && i == f.nIn-1 {
			t = t.Elem()
		}
		if !goodArgType(t) {
			return fmt.Errorf("invalid argument #%d type %s, allowed: number, string, bool", i, t)
		}
	}

	return nil
}

func goodArgType(t reflect.Type) bool {
	kind := t.Kind()
	return kind == reflect.String || kind == reflect.Bool || numberType(kind)
}

func numberType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func joinCompatible(f *funcInfo) error {
	errSig := errors.New("incompatible function signature, expected func(<number>) (string, error) or func(*sqlf.FragmentContext, <number>) (string, error)")
	if f.nOut != 2 && f.outTypes[1] != errorType {
		return errSig
	}
	switch f.nIn {
	case 1:
		if !numberType(f.inTypes[0].Kind()) {
			return errSig
		}
	case 2:
		if !numberType(f.inTypes[1].Kind()) {
			return errSig
		}
	default:
		return errSig
	}
	ctx := newFragmentContext(&Context{
		funcs:    make(map[string]*funcInfo),
		argStore: make([]any, 0),
	}, &Fragment{})
	_, err := evalCall(ctx, f, []any{math.MaxInt32})
	if err == nil || !errors.Is(err, ErrInvalidIndex) {
		return errors.New("never reports ErrInvalidIndex")
	}
	return nil
}
