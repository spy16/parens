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
func (se SymbolExp) Eval(scope *reflection.Scope) (interface{}, error) {
	return scope.Get(se.Symbol)
}
