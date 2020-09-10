package value

import "reflect"

func IsNil(v Any) bool {
	_, isNilType := v.(Nil)
	return v == nil || isNilType
}

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
