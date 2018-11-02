package reflection

import (
	"fmt"
	"reflect"
)

// Call will execute a callable with given args. If the value bound
// to the name is not a callable, ErrNotCallable will be returned.
func Call(rVal reflect.Value, args ...interface{}) (interface{}, error) {
	if rVal.Kind() != reflect.Func {
		return nil, ErrNotCallable
	}

	rType := rVal.Type()
	if !rType.IsVariadic() && rType.NumIn() != len(args) {
		return nil, ErrInvalidNumberOfArgs
	}

	argVals := []reflect.Value{}
	for i := 0; i < rType.NumIn(); i++ {
		argVal := reflect.ValueOf(args[i])
		actualArgType := argVal.Type()
		expectedArgType := rType.In(i)

		if actualArgType != expectedArgType {
			return nil, fmt.Errorf("invalid argument type: expected=%s, actual=%s", expectedArgType, actualArgType)
		}

		argVals = append(argVals, argVal)
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
