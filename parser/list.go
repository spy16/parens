package parser

import (
	"reflect"

	"github.com/spy16/parens/reflection"
)

// ListExp represents a s-exp list containing s-exps.
type ListExp struct {
	List []SExp
}

// Eval evaluates each s-exp in the list and then evaluates the list itself
// as an s-exp.
func (le ListExp) Eval(scope *reflection.Scope) (interface{}, error) {
	if len(le.List) == 0 {
		return le.List, nil
	}

	val, err := le.List[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	reflectVal := reflect.ValueOf(val)
	if reflectVal.Kind() != reflect.Func {
		return nil, reflection.ErrNotCallable
	}

	if macroFn, ok := val.(MacroFunc); ok {
		return macroFn(scope, le.List[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(le.List); i++ {
		arg, err := le.List[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return reflection.Call(reflectVal, args...)
}

// MacroFunc will recieve un-evaluated list of s-expressions.
type MacroFunc func(scope *reflection.Scope, sexps []SExp) (interface{}, error)
