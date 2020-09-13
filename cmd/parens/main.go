package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

func main() {
	ctx := parens.New()

	_ = repl.New(ctx,
		repl.WithBanner("Welcome to Parens!"),
		repl.WithPrompts(">>", " |"),
	).Loop(context.Background())
}
