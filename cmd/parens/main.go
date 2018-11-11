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
		execString(src, exec)
	} else if len(os.Args) == 2 {
		execFile(exec)
	} else {
		runREPL(ctx, exec)
	}

}

func execString(src string, exec *parens.Interpreter) {
	val, err := exec.Execute(src)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(val)
	return
}

func execFile(exec *parens.Interpreter) {
	_, err := exec.ExecuteFile(os.Args[1])
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	return
}

func runREPL(ctx context.Context, exec *parens.Interpreter) {
	exec.DefaultSource = "<REPL>"
	repl, err := parens.NewREPL(exec)
	if err != nil {
		panic(err)
	}
	repl.Banner = "Welcome to Parens REPL!\nType \"(?)\" for help!"
	repl.Start(ctx)
}
