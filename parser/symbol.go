package parser

import (
	"github.com/spy16/parens/reflection"
)

// SymbolExpr represents a symbol.
type SymbolExpr struct {
	Symbol string
}

// ExpType returns s-expression type name.
func (se SymbolExpr) ExpType() string {
	return "symbol"
}

// Eval returns the symbol name itself.
func (se SymbolExpr) Eval(scope *reflection.Scope) (interface{}, error) {
	return scope.Get(se.Symbol)
}
