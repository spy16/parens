package reader

import (
	"github.com/spy16/parens"
)

var defaultSymTable = map[string]parens.Any{
	"nil":   parens.Nil{},
	"true":  parens.Bool(true),
	"false": parens.Bool(false),
}

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

// WithSymbolReader sets the symbol reader macro to be used by the Reader.
// Builds a parens.Symbol if nil.
func WithSymbolReader(m Macro) Option {
	if m == nil {
		return WithBuiltinSymbolReader(defaultSymTable)
	}
	return func(rd *Reader) {
		rd.symReader = m
	}
}

// WithBuiltinSymbolReader configures the default symbol reader with given
// symbol table.
func WithBuiltinSymbolReader(symTable map[string]parens.Any) Option {
	m := symbolReader(symTable)
	return func(rd *Reader) {
		rd.symReader = m
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithNumReader(nil),
		WithBuiltinSymbolReader(defaultSymTable),
	}, opt...)
}
