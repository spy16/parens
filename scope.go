package parens

import (
	"fmt"

	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// NewScope initializes a new scope with given parent scope. parent
// can be nil.
func NewScope(parent parser.Scope) *Scope {
	return &Scope{
		parent: parent,
		vals:   map[string]reflection.Value{},
	}
}

// Scope manages lifetime of values. Scope can inherit values
// from a parent as well.
type Scope struct {
	parent parser.Scope
	vals   map[string]reflection.Value
}

// Root traverses the entire heirarchy of scopes and returns the topmost
// one (i.e., the one with no parent).
func (sc *Scope) Root() parser.Scope {
	if sc.parent == nil {
		return sc
	}

	return sc.parent.Root()
}

// Bind will bind the value to the given name. If a value already
// exists for the given name, it will be overwritten.
func (sc *Scope) Bind(name string, v interface{}) error {
	val := reflection.NewValue(v)
	sc.vals[name] = val

	return nil
}

// Value returns a pointer to the value bound to the given name. If
// name is not available in this scope, request is delgated to the
// parent. If the name is not found anywhere in the hierarchy, error
// will be returned. Modifying the returned pointer will not modify
// the original value.
func (sc *Scope) Value(name string) (*reflection.Value, error) {
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
