package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spy16/parens/lexer"
)

const sampleProgram = `
(println "Hello")
(+ 1.3 10)
(- 10 -100)
`

func main() {
	lxr := lexer.New(sampleProgram)
	spew.Dump(lxr.Tokens())
}
