package main

import (
	"fmt"
	"time"

	"github.com/spy16/parens"
	"github.com/spy16/parens/value"
)

func main() {
	rootCtx := parens.New()

	st := time.Now()
	res, err := rootCtx.Eval(value.Int64(10))
	doneAt := time.Since(st)
	fmt.Println(res, err, doneAt)
}
