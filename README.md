# Parens

[![GoDoc](https://godoc.org/github.com/spy16/parens?status.svg)](https://godoc.org/github.com/spy16/parens) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/parens)](https://goreportcard.com/report/github.com/spy16/parens) ![Go](https://github.com/spy16/parens/workflows/Go/badge.svg?branch=master)


**DEPRECATED**: *This repository is deprecated in favour much better [slurp](github.com/spy16/slurp) project and will be archived/removed soon.*

Parens is a highly customisable, embeddable LISP toolkit.

## Features

* Highly customizable and powerful reader/parser through a read table (Inspired by Clojure) (See [Reader](#reader))
* Built-in data types: nil, bool, string, number, character, keyword, symbol, list.
* Multiple number formats supported: decimal, octal, hexadecimal, radix and scientific notations.
* Full unicode support. Symbols can include unicode characters (Example: `find-δ`, `π` etc.)
  and `🧠`, `🏃` etc. (yes, smileys too).
* Character Literals with support for:
  1. simple literals  (e.g., `\a` for `a`)
  2. special literals (e.g., `\newline`, `\tab` etc.)
  3. unicode literals (e.g., `\u00A5` for `¥` etc.)
* Easy to extend. See [Extending](#extending).
* A macro system.

> Please note that Parens is _NOT_ an implementation of a particular LISP dialect. It provides
> pieces that can be used to build a LISP dialect or can be used as a scripting layer.

## Usage

What can you use it for?

1. Embedded script engine to provide dynamic behavior without requiring re-compilation
   of your application.
2. Business rule engine by exposing very specific & composable rule functions.
3. To build your own LISP dialect.

> Parens requires Go 1.14 or higher.

## Extending

### Reader

Parens reader is inspired by Clojure reader and uses a _read table_. Reader can be extended
to add new syntactical features by adding _reader macros_ to the _read table_. _Reader Macros_
are implementations of `reader.Macro` function type. All syntax that reader can read are 
implemented using _Reader Macros_. Use `SetMacro()` method of `reader.Reader` to override or 
add a custom reader or dispatch macro.

Reader returned by `reader.New(...)`, is configured to support following forms:

* Numbers:
  * Integers use `int64` Go representation and can be specified using decimal, binary
    hexadecimal or radix notations. (e.g., 123, -123, 0b101011, 0xAF, 2r10100, 8r126 etc.)
  * Floating point numbers use `float64` Go representation and can be specified using
    decimal notation or scientific notation. (e.g.: 3.1412, -1.234, 1e-5, 2e3, 1.5e3 etc.)
  * You can override number reader using `WithNumReader()`. 
* Characters: Characters use `rune` or `uint8` Go representation and can be written in 3 ways:
  * Simple: `\a`, `\λ`, `\β` etc.
  * Special: `\newline`, `\tab` etc.
  * Unicode: `\u1267`
* Boolean: `true` or `false` are converted to `Bool` type.
* Nil: `nil` is represented as a zero-allocation empty struct in Go.
* Keywords: Keywords represent symbolic data and start with `:`. (e.g., `:foo`)
* Symbols: Symbols can be used to name a value and can contain any Unicode symbol.
* Lists: Lists are zero or more forms contained within parenthesis. (e.g., `(1 2 3)`, `(1 [])`).

### Evaluation

Parens uses an `Env` for evaluating forms. A form is first macro-expanded and then analysed
to produce an `Expr` that can be evaluated. 

* Macro-expansion can be customised by setting a custom `Expander` implementation. See `parens.WithExpander()`.
* Syntax analysis can be customised (For example, to add special forms), by setting a custom 
  `Analyzer` implementation. See `parens.WithAnalyzer()`.

![I've just received word that the Emperor has dissolved the MIT computer science program permanently.](https://imgs.xkcd.com/comics/lisp_cycles.png)
