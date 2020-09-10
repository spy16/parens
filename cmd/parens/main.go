package main

import (
	"fmt"

	"github.com/spy16/parens"
)

func main() {
	var ev parens.Evaluator
	ctx := parens.NewContext()
	fmt.Println(ev.Eval(ctx, "hello"))
}
