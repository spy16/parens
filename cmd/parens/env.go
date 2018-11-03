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

func makeEnv() *reflection.Env {
	env := reflection.New()
	env.Bind("parens-version", version)

	env.Bind("?", func() string {
		return help
	})

	env.Bind("tokenize", func(src string) ([]lexer.Token, error) {
		return lexer.New(src).Tokens()
	})

	return env
}
