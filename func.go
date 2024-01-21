package sqlf

import (
	"errors"
	"fmt"
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

	joinTested bool  // whether the function has been tested for #join()s
	joinError  error // error to return when the function is not compatible with #join()
}

// JoinCompatibilityError reports whether the function is compatible with #join().
func (f *funcInfo) JoinCompatibilityError() error {
	if f.joinTested {
		return f.joinError
	}
	f.joinError = joinCompatibility(f)
	f.joinTested = true
	return f.joinError
}

// createValueFuncs turns a FuncMap into a map[string]reflect.Value
func createValueFuncs(funcMap FuncMap) (map[string]*funcInfo, error) {
	m := make(map[string]*funcInfo)
	err := addValueFuncs(m, funcMap)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// addValueFuncs adds to values the functions in funcs, converting them to *funcInfos.
func addValueFuncs(out map[string]*funcInfo, in FuncMap) error {
	for name, fn := range in {
		if !goodName(name) {
			return fmt.Errorf("function name %q is not a valid identifier, only letters and underscore are allowed", name)
		}
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			return fmt.Errorf("value for #%s not a function", name)
		}
		typ := v.Type()

		if _, ok := out[name]; ok {
			return fmt.Errorf("function #%s already defined", name)
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
			return fmt.Errorf("function #%s: %w", name, err)
		}

		out[name] = fun
	}
	return nil
}

// goodName reports whether the function name is a valid identifier.
func goodName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r != '_' && !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// goodFunc reports whether the function or method has the right result signature.
func goodFunc(f *funcInfo) error {
	errInvalidFuncOutput := errors.New("invalid signature, expected func(...) (string, error); func(...) string; func(...);")
	switch f.nOut {
	case 0:
		// ok
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
			return fmt.Errorf("unsupported argument type '%s', allowed: number(int*, uint*, float*), string, bool, *sqlf.FragmentContext(as the first argument only)", t)
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

func joinCompatibility(f *funcInfo) error {
	errSig := errors.New("incompatible function signature, expected func(<number>) (string, error) or func(*sqlf.FragmentContext, <number>) (string, error)")
	if f.nOut != 2 || f.outTypes[1] != errorType {
		return errSig
	}
	switch f.nIn {
	case 1:
		if !numberType(f.inTypes[0].Kind()) {
			return errSig
		}
	case 2:
		if !f.contexArg || !numberType(f.inTypes[1].Kind()) {
			return errSig
		}
	default:
		return errSig
	}
	ctx := newFragmentContext(&Context{
		funcs: make(map[string]*funcInfo),
	}, &Fragment{})
	// #join() Assume that the index starts from 1, so 0 is an invalid index,
	// a compatible function should return ErrInvalidIndex
	_, err := evalCall(ctx, f, []any{0})
	if err == nil || !errors.Is(err, ErrInvalidIndex) {
		return errors.New("never reports ErrInvalidIndex")
	}
	return nil
}
