package stdlib

import (
	"fmt"
	"reflect"

	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// Setq is a macro function exposed as 'setq' in parens.
func Setq(scope *reflection.Scope, sexps []parser.SExp) (interface{}, error) {
	if len(sexps) != 2 {
		return nil, fmt.Errorf("expecting symbol and a value")
	}
	symbol, ok := sexps[0].(parser.SymbolExp)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be a symbol, not '%s'", reflect.TypeOf(sexps[0]).String())
	}

	val, err := sexps[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	scope.Bind(symbol.Symbol, val)

	return val, nil
}
