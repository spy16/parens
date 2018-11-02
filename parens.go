package parens

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// New initializes new parens LISP interpreter with given env.
func New(env *reflection.Env) *Interpreter {
	return &Interpreter{
		Env: env,
	}
}

// Interpreter represents the LISP interpreter instance.
type Interpreter struct {
	*reflection.Env
}

// Execute tokenizes, parses and executes the given LISP code.
func (parens *Interpreter) Execute(src string) (interface{}, error) {
	tokens, err := lexer.New(src).Tokens()
	if err != nil {
		return nil, err
	}

	sexp, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	res, err := sexp.Eval(parens.Env)
	if err != nil {
		return nil, err
	}

	return res, nil
}
