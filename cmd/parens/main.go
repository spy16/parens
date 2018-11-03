package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/spy16/parens"
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/reflection"
	"github.com/spy16/parens/repl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	env := makeEnv()
	env.Bind("exit", func() {
		fmt.Println("Bye!")
		cancel()
	})
	interpreter := parens.New(env)

	repl := repl.New(interpreter, nil)
	repl.SetBanner("Welcome to Parens REPL!\nType \"(exit)\" to exit!")
	repl.Start(ctx, os.Stdout, os.Stderr)
}

func makeEnv() *reflection.Env {
	env := reflection.New()
	env.Bind("parens-version", "1.0.0")
	env.Bind("set", func(name string, v interface{}) {
		env.Bind(name, v)
	})
	env.Bind("^", func(base float64, pow float64) float64 {
		return math.Pow(base, pow)
	})

	env.Bind("≠", func(a, b float64) bool {
		return a != b
	})

	env.Bind("tokenize", func(src string) ([]lexer.Token, error) {
		return lexer.New(src).Tokens()
	})

	return env
}
