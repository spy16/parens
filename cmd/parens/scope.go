package main

import (
	"github.com/spy16/parens/lexer"
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

	return scope
}
