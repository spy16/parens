package parser

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/reflection"
)

// SymbolExp represents a symbol.
type SymbolExp struct {
	Symbol string
	Token  lexer.Token
}

// ExpType returns s-expression type name.
func (se SymbolExp) ExpType() string {
	return "symbol"
}

// Eval returns the symbol name itself.
func (se SymbolExp) Eval(env *reflection.Env) (interface{}, error) {
	return env.Get(se.Symbol)
}
