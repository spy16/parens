package main

import (
	"fmt"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reflection"
)

func main() {
	env := reflection.New()

	// Set some values
	env.Bind("pi", 3.1412)
	env.Bind("parens-version", "parens v0.0.1")

	env.Bind("print", func(msg string) {
		fmt.Println(msg)
	})

	env.Bind(">", func(v1, v2 float64) bool {
		return v1 > v2
	})

	env.Bind("<", func(v1, v2 float64) bool {
		return v1 < v2
	})

	ins := parens.New(env)

	_, err := ins.Execute("(print parens-version)")
	if err != nil {
		panic(err)
	}

	isGreater, err := ins.Execute("(> 2 1)")
	if err != nil {
		panic(err)
	}

	fmt.Printf("(> 2 1) =====> %t\n", isGreater)

	_, err = ins.Execute(`(print "hello")`)

}
