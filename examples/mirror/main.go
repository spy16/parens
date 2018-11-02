package main

import (
	"fmt"

	"github.com/spy16/parens/reflection"
)

func main() {
	env := reflection.New()

	panicOnErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	panicOnErr(env.Bind("sample-float", 3.1412))
	panicOnErr(env.Bind("sample-int", 10))
	panicOnErr(env.Bind("hello", func(blah string) string {
		return "hello " + blah
	}))

	fmt.Println(env.GetFloat("sample-float"))
	fmt.Println(env.GetBool("sample"))
	fmt.Println(env)
}
