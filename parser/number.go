package parser

import (
	"strconv"
)

// NumberExpr represents number s-expression.
type NumberExpr struct {
	numStr string
	number interface{}
}

// Eval for a number returns itself.
func (ne NumberExpr) Eval(scope Scope) (interface{}, error) {
	if ne.number == nil {
		num, err := strconv.ParseFloat(ne.numStr, 64)
		if err != nil {
			return nil, err
		}

		ne.number = num
	}

	return ne.number, nil
}
