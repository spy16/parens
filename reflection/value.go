package reflection

import (
	"reflect"
)

// NewValue creates a wrapper Value type for the given value.
func NewValue(v interface{}) Value {
	val := reflect.ValueOf(v)

	return Value{
		Value: val,
	}
}

// Value represents the value associated with a name in the Env.
type Value struct {
	Value reflect.Value
}
