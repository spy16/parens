package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
	"github.com/spy16/parens/value"
)

func main() {
	globals := map[string]value.Any{
		"pi":      value.Float64(3.1416),
		"version": &value.String{Value: "1.0"},
	}

	ctx := parens.New(parens.WithGlobals(globals))

	_ = repl.New(ctx,
		repl.WithBanner("Welcome to Parens!"),
		repl.WithPrompts(">>", " |"),
	).Loop(context.Background())
}
