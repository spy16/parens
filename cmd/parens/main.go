package main

import (
	"context"
	"log"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

func main() {
	globals := map[string]parens.Any{
		"nil":       parens.Nil{},
		"true":      parens.Bool(true),
		"false":     parens.Bool(false),
		"*version*": parens.String("1.0"),
	}

	env := parens.New(parens.WithGlobals(globals, nil))

	err := repl.New(env,
		repl.WithBanner("Welcome to Parens!"),
		repl.WithPrompts(">>", " |"),
	).Loop(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
