package parser

import (
	"github.com/spy16/parens/lexer"
)

func newSymbolExpr(token *lexer.Token) SymbolExpr {
	return SymbolExpr{
		Symbol: token.Value,
	}
}

// SymbolExpr represents a symbol.
type SymbolExpr struct {
	Symbol string
}

// ExpType returns s-expression type name.
func (se SymbolExpr) ExpType() string {
	return "symbol"
}

// Eval returns the symbol name itself.
func (se SymbolExpr) Eval(scope Scope) (interface{}, error) {
	return scope.Get(se.Symbol)
}

func (se SymbolExpr) String() string {
	return se.Symbol
}
