package reader

import (
	"github.com/spy16/parens"
)

// Option for Reader
type Option func(*Reader)

// WithPredefinedSymbols maps a set of symbols to a set of values globally.
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
		WithPredefinedSymbols(nil),
	}, opt...)
}
