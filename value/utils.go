package value

import "reflect"

// IsNil returns true if value is native go `nil` or `Nil{}`.
func IsNil(v Any) bool {
	if v == nil {
		return true
	}

	_, isNilType := v.(Nil)
	return isNilType
}

// IsTruthy returns true if the value has a logical vale of `true`.
func IsTruthy(v Any) bool {
	if IsNil(v) {
		return false
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Bool {
		return rv.Bool()
	}

	return true
}
