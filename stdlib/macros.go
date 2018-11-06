package stdlib

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// Lambda macro is for defining lambdas. (lambda (params) body)
func Lambda(scope *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
	if len(sexps) < 2 {
		return nil, errors.New("at-least two arguments required")
	}

	paramList, ok := sexps[0].(parser.ListExp)
	if !ok {
		return nil, fmt.Errorf("first argument must be list of symbols, not '%s'", reflect.TypeOf(sexps[0]))
	}

	params := []string{}
	for _, entry := range paramList.List {
		sym, ok := entry.(parser.SymbolExp)
		if !ok {
			return nil, fmt.Errorf("param list must contain symbols, not '%s'", reflect.TypeOf(entry))
		}

		params = append(params, sym.Symbol)
	}

	lambdaFunc := func(args ...interface{}) interface{} {
		if len(params) != len(args) {
			panic(fmt.Errorf("requires %d arguments, got %d", len(params), len(args)))
		}

		localScope := reflection.NewScope(scope)
		for i := range params {
			localScope.Bind(params[i], args[i])
		}

		val, err := Begin(localScope, "", sexps[1:])
		if err != nil {
			panic(err)
		}

		return val
	}

	return lambdaFunc, nil
}

// Begin executes all s-exps one by one and returns the result of last evaluation.
func Begin(scope *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
	var val interface{}
	var err error
	for _, sexp := range sexps {
		val, err = sexp.Eval(scope)
		if err != nil {
			return nil, err
		}

	}
	return val, nil
}

// Let creates a new sub-scope from the global scope and executes all the
// sexps inside the new scope. Once the Let block ends, all the names bound
// will be removed. In other words, Let is a begin with local scope.
func Let(scope *reflection.Scope, name string, sexps []parser.SExp) (interface{}, error) {
	localScope := reflection.NewScope(scope)

	return Begin(localScope, name, sexps)
}

// Conditional is commonly know LISP (cond (test1 act1)...) construct.
// Tests can be any exressions that evaluate to non-nil and non-false
// value.
func Conditional(scope *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
	lists := []parser.ListExp{}
	for _, exp := range sexps {
		listExp, ok := exp.(parser.ListExp)

		if !ok {
			return nil, errors.New("all arguments must be lists")
		}
		if len(listExp.List) != 2 {
			return nil, errors.New("each argument must be of the form (test action)")
		}
		lists = append(lists, listExp)
	}

	for _, list := range lists {
		testResult, err := list.List[0].Eval(scope)
		if err != nil {
			return nil, err
		}

		if testResult == nil {
			continue
		}

		if resultBool, ok := testResult.(bool); ok && resultBool == false {
			continue
		}

		return list.List[1].Eval(scope)
	}

	return nil, nil
}

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
