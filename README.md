![Parens](./parens.png)

# Parens

Parens is a LISP-like scripting layer for `Go` (or `Golang`).

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

Then you can run `REPL` by running `parens` command. Or you can run a lisp
file using `parens <filename>` command.


## Usage

Take a look at `cmd/parens/main.go` for a good example.

### Basic Usage

Following is a simple interpreter setup:

```go
scope := reflection.NewScope(nil)

// optional - only if standard functions are needed
stdlib.WithBuiltins(scope)

// custom bind. use any values!
scope.Bind("message", "Hello World!")
scope.Bind("π", 3.1412)

interpreter := parens.New(scope)
interpreter.Execute(`(println message)`)
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
    - [x] Support for `[]` (vector)
    - [x] Enable all UTF-8 characters in symbols
- [ ] Better `parser` package
    - [x] Support for macro functions
    - [x] Support for vectors `[]`
    - [ ] Better error reporting
- [ ] Better `reflection` package
    - [x] Support for variadic functions
    - [ ] Support for methods
    - [ ] Type promotion/conversion
        - [x] `intX` types to `int64` and `float64`
        - [x] `floatX` types to `int64` and `float64`
        - [x] any values to `interface{}` type
        - [ ] `intX` and `floatX` types to `int8`, `int16`, `float32` and vice versa?
- [ ] Performance Benchmark (optimization ?)
- [x] Scopes
- [x] REPL
- [ ] `Go` code generation?
