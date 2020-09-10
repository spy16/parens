package main

import (
	"fmt"

	"github.com/spy16/parens"
)

func main() {
	ctx := parens.NewContext()
	ev := parens.New()
	fmt.Println(ev.Eval(ctx, "hello"))
}
