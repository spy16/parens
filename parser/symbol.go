package parser

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
