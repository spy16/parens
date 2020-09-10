package reader

import (
	"github.com/spy16/parens/value"
)

// Option for Reader
type Option func(*Reader)

// MacroTable is a lookup table that maps the first rune of a macro expression to its
// implementation.
type MacroTable map[rune]Macro

// WithMacros sets the macro table for the reader
func WithMacros(t MacroTable) Option {
	if t == nil {
		t = MacroTable{
			'"':  readString,
			';':  readComment,
			':':  readKeyword,
			'\\': readCharacter,
			'(':  readList,
			')':  UnmatchedDelimiter(),
			'\'': quoteFormReader("quote"),
			'~':  quoteFormReader("unquote"),
			'`':  quoteFormReader("syntax-quote"),
		}
	}

	return func(r *Reader) {
		r.macros = t
	}
}

// WithDispatch sets the dispatch table for the reader
func WithDispatch(t MacroTable) Option {
	if t == nil {
		t = MacroTable{}
	}

	return func(r *Reader) {
		r.dispatch = t
	}
}

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
		WithMacros(nil),
		WithDispatch(nil),
		WithPredefinedSymbols(nil),
	}, opt...)
}
