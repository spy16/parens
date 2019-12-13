# Parens

[![GoDoc](https://godoc.org/github.com/spy16/parens?status.svg)](https://godoc.org/github.com/spy16/parens) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/parens)](https://goreportcard.com/report/github.com/spy16/parens)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fspy16%2Fparens.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fspy16%2Fparens?ref=badge_shield)

Parens is a LISP-like scripting layer for `Go` (or `Golang`).

See [benchmarks](#benchmarks).

## Installation

To embed parens in your application import `github.com/spy16/parens`.

For stand-alone usage, install the parens binary using:

```bash
go get -u -v github.com/spy16/parens/cmd/parens
```

Then you can

1. Run `REPL` by running `parens` command.
2. Run a lisp file using `parens <filename>` command.
3. Execute a LISP string using `parens -e "(+ 1 2)"`


## Usage

Take a look at `cmd/parens/main.go` for a good example.

Check out `examples/` for supported constructs.

## Goals:

### 1. Simple

Should have absolute bare minimum functionality.

```go
scope := parens.NewScope(nil)
parens.ExecuteStr("10")
```

Above snippet gives you an interpreter that understands literals, `(eval expr)`
and `(load <file>)`.

### 2. Flexible

Should be possible to control what is available. Standard functions should be registered
not built-in.

```go
stdlib.RegisterAll(scope)
```

Adding this one line into the previous snippet allows you to include some minimal set
of standard functions like `+`, `-`, `*`, `/` etc. and macros like `let`, `cond`, `do`
etc.

The type of `scope` argument in any `parens` function is the following interface:

```go
// Scope is responsible for managing bindings.
type Scope interface {
	Get(name string) (interface{}, error)
	Bind(name string, v interface{}, doc ...string) error
	Root() Scope
}
```


### 3. Interoperable

Should be able to expose Go values inside LISP and vice versa without custom signatures.

```go
// any Go value can be exposed to interpreter as shown here:
scope.Bind("π", 3.1412)
scope.Bind("banner", "Hello from Parens!")

// functions can be bound directly.
// variadic functions are supported too.
scope.Bind("println", fmt.Println)
scope.Bind("printf", fmt.Printf)
scope.Bind("sin", math.Sin)

// once bound you can use them just as easily:
parens.ExecuteStr("(println banner)", scope)
parens.ExecuteStr(`(printf "value of π is = %f" π)`, scope)
```


### 4. Extensible Semantics

Special constructs like `do`, `cond`, `if` etc. can be added using Macros.

```go
// This is simple implementation of '(do expr*)' special-form from Clojure!
func doMacro(scope parens.Scope, exps []parens.Expr) (interface{}, error) {
    var val interface{}
    var err error
    for _, exp := range exps {
        val, err = exp.Eval(scope)
        if err != nil {
            return nil, err
        }
    }

    return val, nil
}

// register the macro func in the scope.
scope.Bind("do", parens.MacroFunc(doMacro))

// finally use it!
src := `
(do
    (println "Hello from parens")
    (label π 3.1412))
`
// value of 'val' after below statement should be 3.1412
val, _ := parens.ExecuteStr(src, scope)

```

See `stdlib/macros.go` for some built-in macros.

## Parens is *NOT*:

1. An implementaion of a particular LISP dialect (like scheme, common-lisp etc.)
2. A new dialect of LISP


## Benchmarks

| Name                                             | Runs       | Time       | Memory   | Allocations  |
| ------------------------------------------------ | ---------- | ---------- | -------- | ------------ |
| BenchmarkParens_Execute/Execute-8                | 1000000    | 1569 ns/op | 448 B/op | 21 allocs/op |
| BenchmarkParens_Execute/ExecuteExpr-8            | 5000000    | 329 ns/op  | 112 B/op | 5 allocs/op  |
| BenchmarkParens_FunctionCall/DirectCall-8        | 2000000000 | 0.29 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkParens_FunctionCall/CallThroughParens-8 | 2000000    | 827 ns/op  | 224 B/op | 9 allocs/op  |
| BenchmarkNonVariadicCall/Normal-8                | 2000000000 | 0.28 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkNonVariadicCall/Reflection-8            | 5000000    | 370 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkNonVariadicCall/WithTypeConversion-8    | 5000000    | 367 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkVariadicCall/Normal-8                   | 500000000  | 3.22 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkVariadicCall/Reflection-8               | 5000000    | 366 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkVariadicCall/WithTypeConversion-8       | 5000000    | 347 ns/op  | 104 B/op | 4 allocs/op  |

See `v0.0.6` or lower for old benchmarks.

## TODO

- [ ] Better error reporting
- [ ] Optimization
- [ ] `Go` code generation?


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fspy16%2Fparens.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fspy16%2Fparens?ref=badge_large)
