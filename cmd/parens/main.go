package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reflection"
	"github.com/spy16/parens/repl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	env := reflection.New()
	env.Bind("parens-version", "1.0.0")
	env.Bind("exit", func() {
		fmt.Println("Bye!")
		cancel()
	})
	env.Bind("set", func(name string, v interface{}) {
		env.Bind(name, v)
	})
	interpreter := parens.New(env)

	repl := repl.New(interpreter, nil)
	repl.SetBanner("Welcome to Parens REPL!\nType \"(exit)\" to exit!")
	repl.Start(ctx, os.Stdout, os.Stderr)
}
