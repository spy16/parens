package reflection

import (
	"fmt"
	"reflect"
)

// Call will execute a callable with given args. If the value bound
// to the name is not a callable, ErrNotCallable will be returned.
func Call(callable interface{}, args ...interface{}) (interface{}, error) {
	rVal := reflect.ValueOf(callable)
	if rVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("value of kind '%s' is not callable", rVal.Kind())
	}
	rType := rVal.Type()

	argVals, err := makeArgs(rType, args...)
	if err != nil {
		return nil, err
	}

	retVals := rVal.Call(argVals)

	if rType.NumOut() == 0 {
		return nil, nil
	} else if rType.NumOut() == 1 {
		return retVals[0].Interface(), nil
	}

	wrappedRetVals := []interface{}{}
	for _, retVal := range retVals {
		wrappedRetVals = append(wrappedRetVals, retVal.Interface())
	}
	return wrappedRetVals, nil
}

func makeArgs(rType reflect.Type, args ...interface{}) ([]reflect.Value, error) {
	argVals := []reflect.Value{}

	if rType.IsVariadic() {
		nonVariadicLength := rType.NumIn() - 1
		for i := 0; i < nonVariadicLength; i++ {
			convertedArgVal, err := convertValueType(args[i], rType.In(i))
			if err != nil {
				return nil, err
			}

			argVals = append(argVals, convertedArgVal)
		}

		variadicType := rType.In(nonVariadicLength).Elem()
		for i := nonVariadicLength; i < len(args); i++ {
			convertedArgVal, err := convertValueType(args[i], variadicType)
			if err != nil {
				return nil, err
			}

			argVals = append(argVals, convertedArgVal)
		}

		return argVals, nil
	}

	if rType.NumIn() != len(args) {
		return nil, fmt.Errorf("call requires exactly %d arguments, got %d", rType.NumIn(), len(args))
	}

	for i := 0; i < rType.NumIn(); i++ {
		convertedArgVal, err := convertValueType(args[i], rType.In(i))
		if err != nil {
			return nil, err
		}

		argVals = append(argVals, convertedArgVal)
	}

	return argVals, nil
}

func convertValueType(v interface{}, expected reflect.Type) (reflect.Value, error) {
	val := NewValue(v)
	if val.RVal.Type() == expected {
		return val.RVal, nil
	}

	converted, err := val.To(expected.Kind())
	if err != nil {
		if err == ErrConversionImpossible {
			return reflect.Value{}, fmt.Errorf("invalid argument type: expected=%s, actual=%s", expected, val.RVal.Type())
		}
		return reflect.Value{}, err
	}

	return reflect.ValueOf(converted), nil
}
