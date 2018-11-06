package reflection

import (
	"reflect"
)

// NewValue creates a reflection wrapper around given value.
func NewValue(v interface{}) Value {
	return Value{
		RVal: reflect.ValueOf(v),
	}
}

// Value represents every value in parens.
type Value struct {
	RVal reflect.Value
}

// To converts the value to requested kind if possible.
func (val *Value) To(kind reflect.Kind) (interface{}, error) {
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
func (val *Value) ToInt64() (int64, error) {
	if val.isInt() {
		return val.RVal.Int(), nil
	} else if val.isFloat() {
		return int64(val.RVal.Float()), nil
	}

	return 0, ErrConversionImpossible
}

// ToFloat64 attempts converting the value to float64.
func (val *Value) ToFloat64() (float64, error) {
	if val.isFloat() {
		return val.RVal.Float(), nil
	} else if val.isInt() {
		return float64(val.RVal.Int()), nil
	}

	return 0, ErrConversionImpossible
}

// ToBool attempts converting the value to bool.
func (val *Value) ToBool() (bool, error) {
	if isKind(val.RVal, reflect.Bool) {
		return val.RVal.Bool(), nil
	}

	return false, ErrConversionImpossible
}

// ToString attempts converting the value to bool.
func (val *Value) ToString() (string, error) {
	if isKind(val.RVal, reflect.String) {
		return val.RVal.String(), nil
	}

	return "", ErrConversionImpossible
}

func (val *Value) isInt() bool {
	return isKind(val.RVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func (val *Value) isFloat() bool {
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
