package reflection

import "fmt"

// NewScope initializes a new scope with given parent scope. parent
// can be nil.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent: parent,
		vals:   map[string]Value{},
	}
}

// Scope manages lifetime of values. Scope can inherit values
// from a parent as well.
type Scope struct {
	parent *Scope
	vals   map[string]Value
}

// Bind will bind the value to the given name. If a value already
// exists for the given name, it will be overwritten.
func (sc *Scope) Bind(name string, v interface{}) {
	val := NewValue(v)
	sc.vals[name] = val
}

// Value returns a pointer to the value bound to the given name. If
// name is not available in this scope, request is delgated to the
// parent. If the name is not found anywhere in the hierarchy, error
// will be returned. Modifying the returned pointer will not modify
// the original value.
func (sc *Scope) Value(name string) (*Value, error) {
	val, found := sc.vals[name]
	if found {
		return &val, nil
	}

	if sc.parent != nil {
		return sc.parent.Value(name)
	}

	return nil, fmt.Errorf("name '%s' not found", name)
}

// Get returns the actual Go value bound to the given name.
func (sc *Scope) Get(name string) (interface{}, error) {
	val, err := sc.Value(name)
	if err != nil {
		return nil, err
	}

	return val.RVal.Interface(), nil
}

func (sc *Scope) String() string {
	return fmt.Sprintf("Env[size=%d]", len(sc.vals))
}
