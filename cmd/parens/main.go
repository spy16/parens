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

	scope := makeGlobalScope()
	scope.Bind("exit", cancel)
	interpreter := parens.New(scope)

	if len(os.Args) == 2 {
		_, err := interpreter.ExecuteFile(os.Args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		return
	}

	interpreter.DefaultSource = "<REPL>"
	repl := repl.New(interpreter)
	repl.Banner = "Welcome to Parens REPL!\nType \"(?)\" for help!"
	repl.Start(ctx)
}
