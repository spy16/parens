> WIP

# Parens

Parens is a simple lisp implementation in `Go` (or `Golang`). 


## Goals:

1. Simple   - Should have absolute bare minimum functionality
2. Flexible - Should be possible to control what is available
3. Interoperable - Should be able to expose Go values inside LISP and vice versa.

Parens is not:

1. Implementaion of particular LISP dialect (like scheme, common-lisp etc.)
2. A new dialect

## TODO

- [x] Basic working prototype
    - [x] basic `lexer` and `parser`  
    - [x] basic `reflection` package with support for interop b/w `Go` and `lisp` 
- [ ] Better `lexer` package 
    - [ ] Consider whitespaes as delimiters and generate errors
    - [ ] Support for single-quote literals
- [ ] Better `parser` package
    - [ ] Support for macro functions
    - [ ] Better error reporting
- [ ] Better `reflection` package
    - [ ] Type promotion/conversion 
    - [ ] Support for macro functions
- [ ] REPL
- [ ] `Go` code generation?
