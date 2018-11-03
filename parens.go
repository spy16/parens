package parens

import (
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// New initializes new parens LISP interpreter with given env.
func New(env *reflection.Env) *Interpreter {
	return &Interpreter{
		Env:   env,
		Parse: parser.Parse,
	}
}

// ParseFn is responsible for tokenizing and building SExp
// out of tokens.
type ParseFn func(src string) (parser.SExp, error)

// Interpreter represents the LISP interpreter instance. You can
// provide your own implementations of LexFn and ParseFn to extend
// the interpreter.
type Interpreter struct {
	*reflection.Env

	// Parse is used to build SExp/AST from source.
	Parse ParseFn
}

// Execute tokenizes, parses and executes the given LISP code.
func (parens *Interpreter) Execute(src string) (interface{}, error) {
	sexp, err := parens.Parse(src)
	if err != nil {
		return nil, err
	}

	res, err := sexp.Eval(parens.Env)
	if err != nil {
		return nil, err
	}

	return res, nil
}
