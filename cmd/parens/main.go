package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spy16/parens"
)

func main() {
	var src string
	flag.StringVar(&src, "e", "", "Execute source passed in as argument")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	scope := makeGlobalScope()
	scope.Bind("exit", cancel)

	exec := parens.New(scope)
	if len(strings.TrimSpace(src)) > 0 {
		val, err := exec.Execute(src)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(val)
		return
	}

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
