package reader

import (
	"github.com/spy16/parens/value"
)

// Option for Reader
type Option func(*Reader)

// WithPredefinedSymbols maps a set of symbols to a set of values globally.
func WithPredefinedSymbols(ss map[string]value.Any) Option {
	if ss == nil {
		ss = map[string]value.Any{
			"nil":   value.Nil{},
			"true":  value.Bool(true),
			"false": value.Bool(false),
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
