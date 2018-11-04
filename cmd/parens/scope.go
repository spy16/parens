package main

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
	"github.com/spy16/parens/stdlib"
)

const version = "1.0.0"

const help = `
Welcome to Parens!

Type (exit) or Ctrl+D or Ctrl+C to exit the REPL.

See "cmd/parens/main.go" in the github repository for
more information.

https://github.com/spy16/parens
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

	scope.Bind("setq", parser.MacroFunc(stdlib.Setq))

	return scope
}
