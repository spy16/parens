package parens

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// New initializes new parens LISP interpreter with given env.
func New(env *reflection.Env) *Interpreter {
	return &Interpreter{
		Env:   env,
		Lex:   defaultLexFn,
		Parse: defaultParseFn,
	}
}

// LexFn is responsible for reading the source and producing
// tokens.
type LexFn func(src string) ([]lexer.Token, error)

// ParseFn is responsible for building SExp out of tokens.
type ParseFn func(tokens []lexer.Token) (parser.SExp, error)

// Interpreter represents the LISP interpreter instance. You can
// provide your own implementations of LexFn and ParseFn to extend
// the interpreter.
type Interpreter struct {
	*reflection.Env

	// Lex is used to tokenize the given source.
	Lex LexFn

	// Parse is used to build SExp from tokens.
	Parse ParseFn
}

// Execute tokenizes, parses and executes the given LISP code.
func (parens *Interpreter) Execute(src string) (interface{}, error) {
	tokens, err := parens.Lex(src)
	if err != nil {
		return nil, err
	}

	sexp, err := parens.Parse(tokens)
	if err != nil {
		return nil, err
	}

	res, err := sexp.Eval(parens.Env)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func defaultLexFn(src string) ([]lexer.Token, error) {
	return lexer.New(src).Tokens()
}

func defaultParseFn(tokens []lexer.Token) (parser.SExp, error) {
	return parser.Parse(tokens)
}
