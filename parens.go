package parens

import (
	"io/ioutil"
	"strings"

	"github.com/spy16/parens/parser"
)

// New initializes new parens LISP interpreter with given env.
func New(scope parser.Scope) *Interpreter {
	return &Interpreter{
		Scope:         scope,
		Parse:         parser.Parse,
		DefaultSource: "<string>",
	}
}

// ParseFn is responsible for tokenizing and building Expr out of tokens.
type ParseFn func(name, src string) (parser.Expr, error)

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

func (parens *Interpreter) executeSrc(name, src string) (interface{}, error) {
	src = strings.TrimSpace(src)
	expr, err := parens.Parse(name, src)
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
