package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	scope := makeGlobalScope()
	scope.Bind("exit", cancel)
	interpreter := parens.New(scope)

	repl := repl.New(interpreter)
	repl.Banner = "Welcome to Parens REPL!\nType \"(exit)\" or Ctrl+D to exit!"
	repl.Start(ctx)
}
