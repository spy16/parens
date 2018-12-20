package main

import (
	"bufio"
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

	if len(strings.TrimSpace(src)) > 0 {
		execString(src, scope)
	} else if len(os.Args) == 2 {
		execFile(scope)
	} else {
		runREPL(ctx, scope)
	}

}

func execString(src string, env parens.Scope) {
	val, err := parens.Execute(strings.NewReader(src), env)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(val)
	return
}

func execFile(env parens.Scope) {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	defer fh.Close()
	rd := bufio.NewReader(fh)

	_, err = parens.Execute(rd, env)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	return
}

func runREPL(ctx context.Context, env parens.Scope) {
	repl, err := newREPL(env)
	if err != nil {
		panic(err)
	}
	repl.Banner = "Welcome to Parens REPL!\nType \"(?)\" for help!"
	repl.Start(ctx)
}
