package stdlib

import (
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// Setq is a macro that process forms (setq <symbol> <s-exp>). Setq macro
// binds the value after evaluating s-exp to the symbol.
func Setq(scope *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
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

// Inspect dumps the sexps in a formatted manner.
func Inspect(scope *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
	pp.Println(sexps)
	return nil, nil
}
