package parser

// KeywordExpr represents a keyword literal.
type KeywordExpr struct {
	Keyword string
}

// Eval returns the keyword itself.
func (ke KeywordExpr) Eval(scope Scope) (interface{}, error) {
	return ke.Keyword, nil
}

func (ke KeywordExpr) String() string {
	return ke.Keyword
}
