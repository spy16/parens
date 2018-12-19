package parser

import (
	"fmt"
	"reflect"
	"strings"
)

// SymbolExpr represents a symbol.
type SymbolExpr struct {
	Symbol string
}

// Eval returns the symbol name itself.
func (se SymbolExpr) Eval(scope Scope) (interface{}, error) {
	parts := strings.Split(se.Symbol, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid member access symbol. must be of format <parent>.<member>")
	}

	obj, err := scope.Get(parts[0])
	if err != nil {
		return nil, err
	}

	if len(parts) == 1 {
		return obj, nil
	}

	member := resolveMember(reflect.ValueOf(obj), parts[1])
	if !member.IsValid() {
		return nil, fmt.Errorf("member '%s' not found on '%s'", parts[1], parts[0])
	}

	return member.Interface(), nil
}

func (se SymbolExpr) String() string {
	return se.Symbol
}

func resolveMember(obj reflect.Value, name string) reflect.Value {
	firstMatch := func(fxs ...func(string) reflect.Value) reflect.Value {
		for _, fx := range fxs {
			if val := fx(name); val.IsValid() && val.CanInterface() {
				return val
			}
		}

		return reflect.Value{}
	}

	var funcs []func(string) reflect.Value
	if obj.Kind() == reflect.Ptr {
		funcs = append(funcs,
			obj.Elem().FieldByName,
			obj.MethodByName,
			obj.Elem().MethodByName,
		)
	} else {
		funcs = append(funcs,
			obj.FieldByName,
			obj.MethodByName,
		)
	}

	return firstMatch(funcs...)
}
