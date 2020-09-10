package main

import (
	"context"

	"github.com/spy16/parens"
	"github.com/spy16/parens/repl"
)

const banner = "Welcome to Parens!"

func main() {
	repl.New(parens.NewContext(),
		repl.WithBanner(banner),
		repl.WithPrompts(">>>", "  |"),
	).Loop(context.Background())
}
