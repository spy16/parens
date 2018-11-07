package stdlib

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens"
	"github.com/spy16/parens/parser"
)

var macros = []mapEntry{
	entry("begin", parser.MacroFunc(Begin),
		"Usage: (begin expr1 expr2 ...)",
	),
	entry("setq", parser.MacroFunc(Setq),
		"Usage: (setq <symbol> expr)",
	),
	entry("cond", parser.MacroFunc(Conditional),
		"Usage: (cond (test1 action1) (test2 action2)...)",
	),
	entry("let", parser.MacroFunc(Let),
		"Usage: (let expr1 expr2 ...)",
	),
	entry("inspect", parser.MacroFunc(Inspect),
		"Usage: (inspect expr)",
	),
	entry("lambda", parser.MacroFunc(Lambda),
		"Defines a lambda.",
		"Usage: (lambda (params) body)",
		"where params: a list of symbols",
		"      body  : one or more s-expressions",
	),
	entry("doc", parser.MacroFunc(Doc),
		"Dispalys documentation for given symbol if available.",
		"Usage: (doc <symbol>)",
	),
}

// Doc shows doc string associated with a symbol. If not found, returns a message.
func Doc(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("exactly 1 argument required, got %d", len(exprs))
	}

	sym, ok := exprs[0].(parser.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("argument must be a Symbol, not '%s'", reflect.TypeOf(exprs[0]))
	}

	val, err := sym.Eval(scope)
	if err != nil {
		return nil, err
	}

	docStr := scope.Doc(sym.Symbol)
	if len(strings.TrimSpace(docStr)) == 0 {
		docStr = fmt.Sprintf("No documentation available for '%s'", sym.Symbol)
	}

	docStr = fmt.Sprintf("%s\n\nGo Type: %s", docStr, reflect.TypeOf(val))
	return docStr, nil
}

// Lambda macro is for defining lambdas. (lambda (params) body)
func Lambda(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) < 2 {
		return nil, errors.New("at-least two arguments required")
	}

	paramList, ok := exprs[0].(parser.ListExpr)
	if !ok {
		return nil, fmt.Errorf("first argument must be list of symbols, not '%s'", reflect.TypeOf(exprs[0]))
	}

	params := []string{}
	for _, entry := range paramList.List {
		sym, ok := entry.(parser.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("param list must contain symbols, not '%s'", reflect.TypeOf(entry))
		}

		params = append(params, sym.Symbol)
	}

	lambdaFunc := func(args ...interface{}) interface{} {
		if len(params) != len(args) {
			panic(fmt.Errorf("requires %d arguments, got %d", len(params), len(args)))
		}

		localScope := parens.NewScope(scope)
		for i := range params {
			localScope.Bind(params[i], args[i])
		}

		val, err := Begin(localScope, "", exprs[1:])
		if err != nil {
			panic(err)
		}

		return val
	}

	return lambdaFunc, nil
}

// Begin executes all s-exps one by one and returns the result of last evaluation.
func Begin(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	var val interface{}
	var err error
	for _, expr := range exprs {
		val, err = expr.Eval(scope)
		if err != nil {
			return nil, err
		}

	}
	return val, nil
}

// Let creates a new sub-scope from the global scope and executes all the
// exprs inside the new scope. Once the Let block ends, all the names bound
// will be removed. In other words, Let is a begin with local scope.
func Let(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	localScope := parens.NewScope(scope)

	return Begin(localScope, name, exprs)
}

// Conditional is commonly know LISP (cond (test1 act1)...) construct.
// Tests can be any exressions that evaluate to non-nil and non-false
// value.
func Conditional(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	lists := []parser.ListExpr{}
	for _, exp := range exprs {
		listExp, ok := exp.(parser.ListExpr)

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
func Setq(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("expecting symbol and a value")
	}
	symbol, ok := exprs[0].(parser.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be a symbol, not '%s'", reflect.TypeOf(exprs[0]).String())
	}

	val, err := exprs[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	scope.Bind(symbol.Symbol, val)

	return val, nil
}

// Inspect dumps the exprs in a formatted manner.
func Inspect(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	pp.Println(exprs)
	return nil, nil
}
