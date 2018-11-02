> WIP

# Parens

Parens is a simple lisp implementation in `Go` (or `Golang`).
More appropriately, Parens is just a collection of packages that
you can embed in your Golang applications and build a domain language
for your usecase.

## Goals:

1. Simple
    - Should have absolute bare minimum functionality.
2. Flexible
    - Should be possible to control what is available
    - Standard functions should be registered not built-in.
3. Interoperable
    - Should be able to expose Go values inside LISP and vice versa without custom signatures.


Parens is *NOT*:

1. An implementaion of a particular LISP dialect (like scheme, common-lisp etc.)
2. A new dialect if LISP


## Usage

Take a look at `cmd/parens/main.go` for a good example.

Following is a simple interpreter setup:

```go
env := reflection.New()
env.Bind("parens-version", "1.0.0")
env.Bind("print", func(msg string) {
	fmt.Println(msg)
})
interpreter := parens.New(env)
interpreter.Execute("(print parens-version)")
```

## TODO

- [x] Basic working prototype
    - [x] basic `lexer` and `parser`
    - [x] basic `reflection` package with support for interop b/w `Go` and `lisp`
- [x] Better `lexer` package
    - [x] Good unit-test coverage [`97.9%` coverage]
    - [x] Consider whitespaes as delimiters and generate errors
    - [x] Support for `[]` (vector) and `'` (quote)
- [ ] Better `parser` package
    - [ ] Support for macro functions
    - [ ] Better error reporting
    - [ ] Support for single-quote literals
    - [ ] Support for vectors `[]`
- [ ] Better `reflection` package
    - [ ] Type promotion/conversion
    - [ ] Support for macro functions
    - [ ] Performance optimization ?
- [x] REPL
- [ ] `Go` code generation?
