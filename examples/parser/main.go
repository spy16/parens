package main

import (
	"fmt"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

var program = `(print join "hello "world")`

func main() {
	tokens, err := lexer.New(program).Tokens()
	if err != nil {
		panic(err)
	}

	sexp, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	pp.Println(sexp)

	env := reflection.New()
	env.Bind("print", func(msg string) {
		fmt.Println(msg)
	})

	env.Bind("join", func(arg1, arg2 string) string {
		return strings.Join([]string{arg1, arg2}, " ")
	})

	pp.Println(sexp.Eval(env))

}
