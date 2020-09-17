package reader

import (
	"github.com/spy16/parens"
)

// Option values can be used with New() to configure the reader during init.
type Option func(*Reader)

// WithNumReader sets the number reader macro to be used by the Reader. Uses
// the default number reader if nil.
func WithNumReader(m Macro) Option {
	if m == nil {
		m = readNumber
	}
	return func(rd *Reader) {
		rd.numReader = m
	}
}

// WithSymbolFactory sets the symbol factory function for the Reader.
// This function will be used to transform a native Go string into a symbol
// type that satisfies parens.Any.  Builds a parens.Symbol if nil.
func WithSymbolFactory(f func(string) (parens.Any, error)) Option {
	if f == nil {
		f = func(s string) (parens.Any, error) { return parens.Symbol(s), nil }
	}

	return func(rd *Reader) {
		rd.newSymbol = f
	}
}

// WithPredefinedSymbols maps a set of symbols to a set of values globally.
// Reader directly returns the value in the map instead of returning the
// symbol itself.
func WithPredefinedSymbols(ss map[string]parens.Any) Option {
	if ss == nil {
		ss = map[string]parens.Any{
			"nil":   parens.Nil{},
			"true":  parens.Bool(true),
			"false": parens.Bool(false),
		}
	}

	return func(r *Reader) {
		r.predef = ss
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithNumReader(nil),
		WithSymbolFactory(nil),
		WithPredefinedSymbols(nil),
	}, opt...)
}
