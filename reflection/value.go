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

// ToInt attempts converting the value to int64.
func (val *Value) ToInt() (int64, error) {
	if isKind(val.RVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
		return val.RVal.Int(), nil
	}
	return 0, ErrConversionImpossible
}

// ToFloat attempts converting the value to float64.
func (val *Value) ToFloat() (float64, error) {
	if isKind(val.RVal, reflect.Float32, reflect.Float64) {
		return val.RVal.Float(), nil
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

func isKind(rval reflect.Value, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if rval.Kind() == kind {
			return true
		}
	}
	return false
}
