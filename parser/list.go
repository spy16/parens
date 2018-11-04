package parser

import (
	"reflect"

	"github.com/spy16/parens/reflection"
)

// ListExp represents a list (i.e., a function call) s-exp.
type ListExp struct {
	list []SExp
}

// Eval evaluates each s-exp in the list and then evaluates the list itself
// as an s-exp.
func (le ListExp) Eval(scope *reflection.Scope) (interface{}, error) {
	if len(le.list) == 0 {
		return le.list, nil
	}

	val, err := le.list[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	reflectVal := reflect.ValueOf(val)
	if reflectVal.Kind() != reflect.Func {
		return nil, reflection.ErrNotCallable
	}

	if macroFn, ok := val.(MacroFunc); ok {
		return macroFn(scope, le.list[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(le.list); i++ {
		arg, err := le.list[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return reflection.Call(reflectVal, args...)
}

// MacroFunc will recieve un-evaluated list of s-expressions.
type MacroFunc func(scope *reflection.Scope, sexps []SExp) (interface{}, error)
