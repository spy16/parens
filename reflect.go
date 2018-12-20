package parens

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrNameNotFound is returned when a lookup is performed with a
	// non-bound name.
	ErrNameNotFound = errors.New("name not bound to a value")

	// ErrNotCallable is returned when a Call is attempted on a non-
	// callable value.
	ErrNotCallable = errors.New("value is not callable")

	// ErrConversionImpossible is returned when the Value type cannot be
	// converted to the expected type.
	ErrConversionImpossible = errors.New("cannot be converted")

	// ErrInvalidNumberOfArgs is returned when a function call is attempted
	// with invalid number of arguments.
	ErrInvalidNumberOfArgs = errors.New("invalid number of arguments")
)

// reflectCall will execute a callable with given args. If the value bound
// to the name is not a callable, ErrNotCallable will be returned.
func reflectCall(callable interface{}, args ...interface{}) (interface{}, error) {
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
	val := newValue(v)
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

func newValue(v interface{}) reflectVal {
	return reflectVal{
		RVal: reflect.ValueOf(v),
	}
}

type reflectVal struct {
	RVal reflect.Value
}

// To converts the value to requested kind if possible.
func (val *reflectVal) To(kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.Int, reflect.Int64:
		return val.ToInt64()
	case reflect.Float64:
		return val.ToFloat64()
	case reflect.String:
		return val.ToString()
	case reflect.Bool:
		return val.ToBool()
	case reflect.Interface:
		return val.RVal.Interface(), nil
	default:
		return nil, ErrConversionImpossible
	}
}

// ToInt64 attempts converting the value to int64.
func (val *reflectVal) ToInt64() (int64, error) {
	if val.isInt() {
		return val.RVal.Int(), nil
	} else if val.isFloat() {
		return int64(val.RVal.Float()), nil
	}

	return 0, ErrConversionImpossible
}

// ToFloat64 attempts converting the value to float64.
func (val *reflectVal) ToFloat64() (float64, error) {
	if val.isFloat() {
		return val.RVal.Float(), nil
	} else if val.isInt() {
		return float64(val.RVal.Int()), nil
	}

	return 0, ErrConversionImpossible
}

// ToBool attempts converting the value to bool.
func (val *reflectVal) ToBool() (bool, error) {
	if isKind(val.RVal, reflect.Bool) {
		return val.RVal.Bool(), nil
	}

	return false, ErrConversionImpossible
}

// ToString attempts converting the value to bool.
func (val *reflectVal) ToString() (string, error) {
	if isKind(val.RVal, reflect.String) {
		return val.RVal.String(), nil
	}

	return "", ErrConversionImpossible
}

func (val *reflectVal) isInt() bool {
	return isKind(val.RVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func (val *reflectVal) isFloat() bool {
	return isKind(val.RVal, reflect.Float32, reflect.Float64)
}

func isKind(rval reflect.Value, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if rval.Kind() == kind {
			return true
		}
	}
	return false
}
