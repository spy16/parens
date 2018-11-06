package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spy16/parens"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	scope := makeGlobalScope()
	scope.Bind("exit", cancel)

	exec := parens.New(scope)
	if len(os.Args) == 2 {
		_, err := exec.ExecuteFile(os.Args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		return
	}

	exec.DefaultSource = "<REPL>"
	repl := parens.NewREPL(exec)
	repl.Banner = "Welcome to Parens REPL!\nType \"(?)\" for help!"
	repl.Start(ctx)
}
