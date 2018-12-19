package parens

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/spy16/parens/parser"
)

// New initializes new parens LISP interpreter with given env.
func New(scope parser.Scope) *Interpreter {
	exec := &Interpreter{
		Scope:         scope,
		Parse:         parser.ParseModule,
		DefaultSource: "<string>",
	}

	loadFile := func(file string) interface{} {
		val, err := exec.ExecuteFile(file)
		if err != nil {
			panic(err)
		}

		return val
	}

	evalStr := func(val interface{}) interface{} {
		expr, ok := val.(parser.Expr)
		if !ok {
			return val
		}

		res, err := expr.Eval(exec.Scope)
		if err != nil {
			panic(err)
		}

		return res
	}

	scope.Bind("load", loadFile,
		"Reads and executes the file in the current scope",
		"Example: (load \"sample.lisp\")",
	)

	scope.Bind("eval", evalStr,
		"Executes given LISP string in the current scope",
		"Usage: (eval <form>)",
	)
	return exec
}

// ParseFn is responsible for tokenizing and building Expr out of tokens.
type ParseFn func(name string, src io.RuneScanner) (parser.Expr, error)

// Interpreter represents the LISP interpreter instance. You can provide
// your own implementations of ParseFn to extend the interpreter.
type Interpreter struct {
	Scope         parser.Scope
	Parse         ParseFn
	DefaultSource string
}

// Execute tokenizes, parses and executes the given LISP code.
func (parens *Interpreter) Execute(src string) (interface{}, error) {
	return parens.executeSrc(parens.DefaultSource, src)
}

// ExecuteFile reads, tokenizes, parses and executes the contents of the given file.
func (parens *Interpreter) ExecuteFile(file string) (interface{}, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return parens.executeSrc(file, string(data))
}

// ExecuteExpr executes the given expr using the appropriate scope.
func (parens *Interpreter) ExecuteExpr(expr parser.Expr) (interface{}, error) {
	return expr.Eval(parens.Scope)
}

func (parens *Interpreter) executeSrc(name, src string) (interface{}, error) {
	src = strings.TrimSpace(src)
	expr, err := parens.Parse(name, strings.NewReader(src))
	if err != nil {
		return nil, err
	}

	var res interface{}
	var evalErr error
	safeWrapper := func() {
		defer func() {
			if v := recover(); v != nil {
				if err, ok := v.(error); ok {
					evalErr = err
				} else {
					evalErr = fmt.Errorf("panic: %v", v)
				}
			}
		}()

		res, err = expr.Eval(parens.Scope)
		evalErr = err
	}

	safeWrapper()
	if evalErr != nil {
		return nil, evalErr
	}

	return res, nil
}
