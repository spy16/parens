package stdlib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens"
)

var core = []mapEntry{
	// logical constants
	entry("true", true,
		"Represents logical true",
	),
	entry("false", false,
		"Represents logical false",
	),
	entry("nil", false,
		"Represents logical false. Same as false",
	),

	// core macros
	entry("do", parens.MacroFunc(Do),
		"Usage: (do expr1 expr2 ...)",
	),
	entry("quote", parens.MacroFunc(Quote),
		"Usage: (quote expr)",
	),
	entry("label", parens.MacroFunc(Label),
		"Usage: (label <symbol> expr)",
	),
	entry("global", parens.MacroFunc(Global),
		"Usage: (global <symbol> expr)",
	),
	entry("cond", parens.MacroFunc(Conditional),
		"Usage: (cond (test1 action1) (test2 action2)...)",
	),
	entry("let", parens.MacroFunc(Let),
		"Usage: (let expr1 expr2 ...)",
	),
	entry("inspect", parens.MacroFunc(Inspect),
		"Usage: (inspect expr)",
	),
	entry("lambda", parens.MacroFunc(Lambda),
		"Defines a lambda.",
		"Usage: (lambda (params) body)",
		"where params: a list of symbols",
		"      body  : one or more s-expressions",
	),
	entry("defn", parens.MacroFunc(Defn),
		"Defines a named function",
		"Usage: (defn <name> [params] body)",
	),
	entry("doc", parens.MacroFunc(Doc),
		"Displays documentation for given symbol if available.",
		"Usage: (doc <symbol>)",
	),
	entry("dump-scope", parens.MacroFunc(dumpScope),
		"Formats and displays the entire scope",
	),
	entry("->", parens.MacroFunc(ThreadFirst)),
	entry("->>", parens.MacroFunc(ThreadLast)),

	// core functions
	entry("type", reflect.TypeOf),
}

// Quote prevents the expr from being executed until unquoted.
func Quote(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("exactly 1 argument required")
	}

	return exprs[0], nil
}

// LoadFile returns function that reads and executes a lisp file.
func LoadFile(env parens.Scope) func(file string) interface{} {
	return func(file string) interface{} {
		fh, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		rd := bufio.NewReader(fh)

		val, err := parens.Execute(rd, env)
		if err != nil {
			panic(err)
		}

		return val
	}
}

// Eval returns the (eval <val>) function.
func Eval(env parens.Scope) func(val interface{}) interface{} {
	return func(val interface{}) interface{} {
		expr, ok := val.(parens.Expr)
		if !ok {
			return val
		}

		res, err := expr.Eval(env)
		if err != nil {
			panic(err)
		}

		return res
	}
}

// ThreadFirst macro appends first evaluation result as first argument of next function
// call.
func ThreadFirst(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	return thread(true, scope, exprs)
}

// ThreadLast macro appends first evaluation result as last argument of next function
// call.
func ThreadLast(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	return thread(false, scope, exprs)
}

func thread(first bool, scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("at-least 1 argument required")
	}

	var result interface{}
	for i := 1; i < len(exprs); i++ {
		lst, ok := exprs[i].(parens.List)
		if !ok {
			return nil, fmt.Errorf("argument %d must be a function call, not '%s'", i, reflect.TypeOf(exprs[i]))
		}

		val, err := exprs[i-1].Eval(scope)
		if err != nil {
			return nil, err
		}
		res := anyExpr{val: val}

		nextCall := parens.List([]parens.Expr{lst[0]})

		if first {
			nextCall = append(nextCall, res)
			nextCall = append(nextCall, lst[1:]...)
		} else {
			nextCall = append(nextCall, lst[1:]...)
			nextCall = append(nextCall, res)
		}

		result, err = nextCall.Eval(scope)
		if err != nil {
			return nil, err
		}
		exprs[i] = anyExpr{val: result}
	}

	return result, nil
}

type anyExpr struct {
	val interface{}
}

func (ae anyExpr) Eval(scope parens.Scope) (interface{}, error) {
	return ae.val, nil
}

// Doc shows doc string associated with a symbol. If not found, returns a message.
func Doc(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("exactly 1 argument required, got %d", len(exprs))
	}

	sym, ok := exprs[0].(parens.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument must be a Symbol, not '%s'", reflect.TypeOf(exprs[0]))
	}

	val, err := sym.Eval(scope)
	if err != nil {
		return nil, err
	}

	var docStr string
	if swd, ok := scope.(scopeWithDoc); ok {
		docStr = swd.Doc(string(sym))
	}
	if len(strings.TrimSpace(docStr)) == 0 {
		docStr = fmt.Sprintf("No documentation available for '%s'", string(sym))
	}

	docStr = fmt.Sprintf("%s\n\nGo Type: %s", docStr, reflect.TypeOf(val))
	return docStr, nil
}

// Defn macro is for defining named functions. It defines a lambda and binds it with
// the given name into the scope.
func Defn(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) < 3 {
		return nil, fmt.Errorf("3 or more arguments required, got %d", len(exprs))
	}

	sym, ok := exprs[0].(parens.Symbol)
	if !ok {
		return nil, fmt.Errorf("first argument must be symbol, not '%s'", reflect.TypeOf(exprs[0]))
	}

	lambda, err := Lambda(scope, exprs[1:])
	if err != nil {
		return nil, err
	}

	scope.Bind(string(sym), lambda)
	return string(sym), nil
}

// Lambda macro is for defining lambdas. (lambda (params) body)
func Lambda(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) < 2 {
		return nil, errors.New("at-least two arguments required")
	}

	paramList, ok := exprs[0].(parens.Vector)
	if !ok {
		return nil, fmt.Errorf("first argument must be list of symbols, not '%s'", reflect.TypeOf(exprs[0]))
	}

	params := []string{}
	for _, entry := range paramList {
		sym, ok := entry.(parens.Symbol)
		if !ok {
			return nil, fmt.Errorf("param list must contain symbols, not '%s'", reflect.TypeOf(entry))
		}

		params = append(params, string(sym))
	}

	lambdaFunc := func(args ...interface{}) interface{} {
		if len(params) != len(args) {
			panic(fmt.Errorf("requires %d arguments, got %d", len(params), len(args)))
		}

		localScope := parens.NewScope(scope)
		for i := range params {
			localScope.Bind(params[i], args[i])
		}

		val, err := Do(localScope, exprs[1:])
		if err != nil {
			panic(err)
		}

		return val
	}

	return lambdaFunc, nil
}

// Do executes all s-exps one by one and returns the result of last evaluation.
func Do(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
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
// will be removed. In other words, Let is a Do with local scope.
func Let(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	localScope := parens.NewScope(scope)

	return Do(localScope, exprs)
}

// Conditional is commonly know LISP (cond (test1 act1)...) construct.
// Tests can be any exressions that evaluate to non-nil and non-false
// value.
func Conditional(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	lists := []parens.List{}
	for _, exp := range exprs {
		listExp, ok := exp.(parens.List)

		if !ok {
			return nil, errors.New("all arguments must be lists")
		}
		if len(listExp) != 2 {
			return nil, errors.New("each argument must be of the form (test action)")
		}
		lists = append(lists, listExp)
	}

	for _, list := range lists {
		testResult, err := list[0].Eval(scope)
		if err != nil {
			return nil, err
		}

		if testResult == nil {
			continue
		}

		if resultBool, ok := testResult.(bool); ok && resultBool == false {
			continue
		}

		return list[1].Eval(scope)
	}

	return nil, nil
}

// Label binds the result of evaluating second argument to the symbol passed in as
// first argument in the current scope.
func Label(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	return labelInScope(scope, exprs)
}

// Global binds the result of evaluating second argument to the symbol passed in as
// first argument in the global scope.
func Global(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	return labelInScope(scope.Root(), exprs)
}

// Inspect dumps the exprs in a formatted manner.
func Inspect(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	pp.Println(exprs)
	return nil, nil
}

func dumpScope(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	return fmt.Sprint(scope), nil
}

func labelInScope(scope parens.Scope, exprs []parens.Expr) (interface{}, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("expecting symbol and a value")
	}
	symbol, ok := exprs[0].(parens.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be a symbol, not '%s'", reflect.TypeOf(exprs[0]).String())
	}

	val, err := exprs[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	scope.Bind(string(symbol), val)

	return val, nil
}

type scopeWithDoc interface {
	Doc(name string) string
}
