package main

import (
	"fmt"

	"github.com/spy16/parens"
)

func main() {
	ev := parens.New()
	fmt.Println(ev.Eval("hello"))
}
