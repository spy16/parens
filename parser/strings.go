package parser

// StringExpr represents single and double quoted strings.
type StringExpr struct {
	Value string
}

// Eval returns unquoted version of the STRING token.
func (se StringExpr) Eval(_ Scope) (interface{}, error) {
	return se.Value, nil
}
