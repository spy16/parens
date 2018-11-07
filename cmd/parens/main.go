package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
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

	prompter := makePrompter(cancel)
	repl.ReadIn = func() (string, error) {
		val := prompter.Input()
		return val, nil
	}

	repl.Start(ctx)
}

func makePrompter(onCtrlC func()) *prompt.Prompt {
	promptExitFunc := func(_ *prompt.Buffer) {
		fmt.Println("Bye!")
		onCtrlC()
		os.Exit(0)
	}

	prompter := prompt.New(nil, func(doc prompt.Document) []prompt.Suggest {
		return nil
	},
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  promptExitFunc,
		}))

	return prompter
}
