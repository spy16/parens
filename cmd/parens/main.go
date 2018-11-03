package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	env := makeEnv()
	env.Bind("exit", cancel)
	interpreter := parens.New(env)

	repl := repl.New(interpreter)
	repl.Banner = "Welcome to Parens REPL!\nType \"(exit)\" to exit!"
	repl.Start(ctx)
}
