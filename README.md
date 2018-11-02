> WIP

# Parens

Parens is a simple lisp implementation in `Go` (or `Golang`).

Parens is not:

1. Implementaion of particular LISP dialect (like scheme, common-lisp etc.)
2. A new dialect

Parens is just a collection of `lexer` and `parser` packages.

## Goals:

1. Simple
    - Should have absolute bare minimum functionality.
2. Flexible
    - Should be possible to control what is available
    - Standard functions should be registered not build-in.
3. Interoperable
    - Should be able to expose Go values inside LISP and vice versa without custom signatures.

## TODO

- [x] Basic working prototype
    - [x] basic `lexer` and `parser`
    - [x] basic `reflection` package with support for interop b/w `Go` and `lisp`
- [x] Better `lexer` package
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
- [x] REPL
- [ ] `Go` code generation?
