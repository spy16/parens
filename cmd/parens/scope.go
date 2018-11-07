package main

import (
	"github.com/spy16/parens"
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/stdlib"
)

const version = "1.0.0"

const help = `
Welcome to Parens!

Type (exit) or Ctrl+C to exit the REPL.

Use (dump-scope) to see the list of symbols available in
the current scope.

Use (doc <symbol>) to get help about symbols in scope.

See "cmd/parens/main.go" in the github repository for
more information.

https://github.com/spy16/parens
`

func makeGlobalScope() *parens.Scope {
	scope := parens.NewScope(nil)
	scope.Bind("parens-version", version)

	scope.Bind("?", func() string {
		return help
	})

	scope.Bind("tokenize", func(src string) ([]lexer.Token, error) {
		return lexer.New(src).Tokens()
	})

	stdlib.RegisterAll(scope)
	return scope
}
