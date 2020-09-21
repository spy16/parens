package main

import (
	"context"
	"log"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

var globals = map[string]parens.Any{
	"nil":       parens.Nil{},
	"true":      parens.Bool(true),
	"false":     parens.Bool(false),
	"*version*": parens.String("1.0"),
}

func main() {
	env := parens.New(parens.WithGlobals(globals, nil))

	r := repl.New(env,
		repl.WithBanner("Welcome to Parens!"),
		repl.WithPrompts(">>", " |"),
	)

	if err := r.Loop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
