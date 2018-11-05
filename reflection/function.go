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
		return nil, ErrNotCallable
	}

	rType := rVal.Type()
	if !rType.IsVariadic() && rType.NumIn() != len(args) {
		return nil, ErrInvalidNumberOfArgs
	}

	argVals := []reflect.Value{}
	for i := 0; i < rType.NumIn(); i++ {
		convertedArgVal, err := convertValueType(args[i], rType.In(i))
		if err != nil {
			return nil, err
		}

		argVals = append(argVals, convertedArgVal)
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

func convertValueType(v interface{}, expected reflect.Type) (reflect.Value, error) {
	val := NewValue(v)
	if val.RVal.Type() == expected {
		return val.RVal, nil
	}

	if expected.Kind() == reflect.Interface {
		return reflect.ValueOf(val.RVal.Interface()), nil
	}

	return reflect.Value{}, fmt.Errorf("invalid argument type: expected=%s, actual=%s", expected, val.RVal.Type())
}
