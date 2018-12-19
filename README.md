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
scope := reflection.NewScope(nil)
exec := parens.New(scope)
exec.Execute("10")
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

The type of `scope` argument in `parens.New(scope)` is the following interface:

```go
// Scope is responsible for managing bindings.
type Scope interface {
	Get(name string) (interface{}, error)
	Doc(name string) string
	Bind(name string, v interface{}, doc ...string) error
	Root() Scope
}
```

So, if you wanted to do some dynamic resolution during `Get(name)` (e.g. If you wanted to return
*method* `Print` of object `stdout` when `Get("stdout.Print")` is called), you can easily implement
this interface and pass it to `parens.New`.


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
exec.Execute("(println banner)")
exec.Execute(`(printf "value of π is = %f" π)`)
```


### 4. Extensible Semantics

Special constructs like `do`, `cond`, `if` etc. can be added using Macros.

```go
// This is simple implementation of '(do expr*)' special-form from Clojure!
func doMacro(scope *reflection.Scope, callName string, exps []parser.Expr) (interface{}, error) {
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
scope.Bind("do", parser.MacroFunc(doMacro))

// finally use it!
src := `
(do
    (println "Hello from parens")
    (label π 3.1412))
`
// value of 'val' after below statement should be 3.1412
val, _ := exec.Execute(src)

```

See `stdlib/macros.go` for some built-in macros.

## Parens is *NOT*:

1. An implementaion of a particular LISP dialect (like scheme, common-lisp etc.)
2. A new dialect of LISP


## Benchmarks

| Name                                             | Runs       | Time       | Memory   | Allocations  |
| ------------------------------------------------ | ---------- | ---------- | -------- | ------------ |
| BenchmarkParens_Execute/Execute-8                | 1000000    | 1977 ns/op | 464 B/op | 21 allocs/op |
| BenchmarkParens_Execute/ExecuteExpr-8            | 5000000    | 323 ns/op  | 112 B/op | 5 allocs/op  |
| BenchmarkParens_FunctionCall/DirectCall-8        | 2000000000 | 0.30 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkParens_FunctionCall/CallThroughParens-8 | 2000000    | 844 ns/op  | 224 B/op | 9 allocs/op  |
| BenchmarkNonVariadicCall/Normal-8                | 2000000000 | 0.28 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkNonVariadicCall/Reflection-8            | 5000000    | 342 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkNonVariadicCall/WithTypeConversio-8     | 5000000    | 343 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkVariadicCall/Normal-8                   | 500000000  | 3.77 ns/op | 0 B/op   | 0 allocs/op  |
| BenchmarkVariadicCall/Reflection-8               | 5000000    | 350 ns/op  | 104 B/op | 4 allocs/op  |
| BenchmarkVariadicCall/WithTypeConversion-8       | 5000000    | 354 ns/op  | 104 B/op | 4 allocs/op  |

See `v0.0.6` or lower for old benchmarks.

## TODO

- [ ] Better way to map error returns from Go functios to LISP
    - [ ] `Result<interface{}, error>` type of design in Rust ?
- [ ] Better `parser` package
    - [x] Support for macro functions
    - [x] Support for vectors `[]`
    - [ ] Better error reporting
- [ ] Better `reflection` package
    - [x] Support for variadic functions
    - [x] Support for methods
    - [ ] Type promotion/conversion
        - [x] `intX` types to `int64` and `float64`
        - [x] `floatX` types to `int64` and `float64`
        - [x] any values to `interface{}` type
        - [ ] `intX` and `floatX` types to `int8`, `int16`, `float32` and vice versa?
- [ ] Optimization
  - [x] Performance Benchmark
  - [x] optimization ?
- [ ] `Go` code generation?


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fspy16%2Fparens.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fspy16%2Fparens?ref=badge_large)