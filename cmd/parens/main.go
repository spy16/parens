package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
	"github.com/spy16/parens/value"
)

func main() {
	globals := map[string]value.Any{
		"nil":       value.Nil{},
		"true":      value.Bool(true),
		"false":     value.Bool(false),
		"*version*": value.String("1.0"),
	}

	env := parens.New(parens.WithGlobals(globals, nil))

	_ = repl.New(env,
		repl.WithBanner("Welcome to Parens!"),
		repl.WithPrompts(">>", " |"),
	).Loop(context.Background())
}
