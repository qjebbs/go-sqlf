package sqlf

import (
	"fmt"
	"reflect"
)

func evalFunction(ctx *FragmentContext, name string, args []any) (string, error) {
	function, ok := ctx.Global.funcs[name]
	if !ok {
		return "", fmt.Errorf("%q is not a defined function", name)
	}
	return evalCall(ctx, function, name, args)
}

func evalCall(ctx *FragmentContext, fun reflect.Value, name string, args []any) (string, error) {
	typ := fun.Type()
	ctxArg := typ.NumIn() > 0 && typ.In(0) == contextPointerType
	if ctxArg {
		args = append([]any{ctx}, args...)
	}
	numIn := len(args)
	numFixed := typ.NumIn()
	if typ.IsVariadic() {
		numFixed-- // last arg is the variadic one.
		if numIn < numFixed {
			return "", fmt.Errorf("wrong number of args for #%s: want at least %d got %d", name, typ.NumIn()-1, len(args))
		}
	} else if numIn != typ.NumIn() {
		return "", fmt.Errorf("wrong number of args for #%s: want %d got %d", name, typ.NumIn(), numIn)
	}

	unwrap := func(v reflect.Value) reflect.Value {
		if v.Type() == reflectValueType {
			v = v.Interface().(reflect.Value)
		}
		return v
	}

	// Build the arg list.
	var err error
	argv := make([]reflect.Value, numIn)
	// Fixed args first.
	i := 0
	for ; i < numFixed && i < len(args); i++ {
		if i == 0 && ctxArg {
			argv[i] = reflect.ValueOf(args[0])
			continue
		}
		inType := typ.In(i)
		argv[i], err = evalArg(inType, args[i])
		if err != nil {
			return "", fmt.Errorf("arg %d has wrong type for #%s: %w", i, name, err)
		}
	}
	// Now the ... args.
	if typ.IsVariadic() {
		inType := typ.In(typ.NumIn() - 1).Elem() // Argument is a slice.
		for ; i < len(args); i++ {
			argv[i], err = evalArg(inType, args[i])
			if err != nil {
				return "", fmt.Errorf("arg %d has wrong type for #%s: %w", i, name, err)
			}
		}
	}
	v, err := safeCall(fun, argv)
	if err != nil {
		return "", fmt.Errorf("error calling #%s: %w", name, err)
	}
	valAny := unwrap(v).Interface()
	val, ok := valAny.(string)
	if !ok {
		return "", fmt.Errorf("function #%s returned %T, not string", name, valAny)
	}
	return val, nil
}

// safeCall runs fun.Call(args), and returns the resulting value and error, if
// any. If the call panics, the panic value is returned as an error.
func safeCall(fun reflect.Value, args []reflect.Value) (val reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	ret := fun.Call(args)
	if len(ret) == 2 && !ret[1].IsNil() {
		return ret[0], ret[1].Interface().(error)
	}
	return ret[0], nil
}

func evalArg(typ reflect.Type, arg any) (reflect.Value, error) {
	var (
		v       reflect.Value
		argType string
	)
	// fmt.Printf("convert arg %T(%v) to %s\n", arg, arg, typ.Name())
	switch arg := arg.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		argType = "number"
		v = reflect.ValueOf(arg).Convert(typ)
	case bool:
		argType = "bool"
		if typ.Kind() == reflect.Bool {
			v = reflect.ValueOf(arg)
		}
	case string:
		argType = "string"
		if typ.Kind() == reflect.String {
			v = reflect.ValueOf(arg)
		}
	case nil:
		argType = "nil"
		if typ.Kind() == reflect.Interface {
			v = reflect.Zero(typ)
		}
	default:
		argType = reflect.TypeOf(arg).Name()
	}
	if v.IsValid() {
		return v, nil
	}
	return v, fmt.Errorf("can't assign %s to %s", argType, typ)
}
