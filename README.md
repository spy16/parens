> WIP

![Parens](./parens.png)

# Parens

Parens is a simple embeddable LISP implementation in `Go` (or `Golang`).
More appropriately, Parens is a scripting layer for `Go` applications which
uses LISP syntax.

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
2. A new dialect of LISP


## Installation

Parens is not meant for stand-alone usage. But there is a REPL which is
meant to showcase features of parens and can be installed as below:

```bash
go get -u -v github.com/spy16/parens/cmd/parens
```


## Usage

Take a look at `cmd/parens/main.go` for a good example.

### Basic Usage

Following is a simple interpreter setup:

```go
scope := reflection.NewScope(nil)
stdlib.WithBuiltins(scope)
interpreter := parens.New(scope)
interpreter.Execute("(print parens-version)")
```

### Macros

Golang functions can be registered as macros into the interpreter
as shown below (for more examples, see `./stdlib/macros.go`):

```go
func inspect(_ *reflection.Scope, _ string, sexps []parser.SExp) (interface{}, error) {
    spew.Dump(sexps)
    return nil, nil
}

scope := reflection.NewScope(nil)
scope.Bind("inspect", parser.MacroFunc(inspect))
```

## TODO

- [x] Basic working prototype
    - [x] basic `lexer` and `parser`
    - [x] basic `reflection` package with support for interop b/w `Go` and `lisp`
- [x] Better `lexer` package
    - [x] Good unit-test coverage [`97.9%` coverage]
    - [x] Consider whitespaes as delimiters and generate errors
    - [x] Support for `[]` (vector) and `'` (quote)
    - [x] Enable all UTF-8 characters in symbols
- [ ] Better `parser` package
    - [x] Support for macro functions
    - [x] Support for vectors `[]`
    - [ ] Better error reporting
    - [ ] Support for single-quote literals
- [ ] Better `reflection` package
    - [ ] Type promotion/conversion
    - [ ] Performance optimization ?
- [x] Scopes
- [x] REPL
- [ ] `Go` code generation?
