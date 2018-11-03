package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spy16/parens"
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
