package main

import (
	"fmt"
	"reflect"

	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

const version = "1.0.0"

const help = `
Welcome to Parens!

Functions:
1. (exit)         - Exit the REPL
2. (tokenize src) - Tokenize src and display
`

func makeGlobalScope() *reflection.Scope {
	scope := reflection.NewScope(nil)
	scope.Bind("parens-version", version)

	scope.Bind("?", func() string {
		return help
	})

	scope.Bind("tokenize", func(src string) ([]lexer.Token, error) {
		return lexer.New(src).Tokens()
	})

	scope.Bind("setq", parser.MacroFunc(func(scope *reflection.Scope, sexps []parser.SExp) (interface{}, error) {
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
	}))

	return scope
}
