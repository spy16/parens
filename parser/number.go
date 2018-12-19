package parser

import (
	"fmt"
	"strconv"
)

// NumberExpr represents number s-expression.
type NumberExpr struct {
	NumStr string
	Number interface{}
}

// Eval for a number returns itself.
func (ne NumberExpr) Eval(scope Scope) (interface{}, error) {
	if ne.Number == nil {
		num, err := strconv.ParseFloat(ne.NumStr, 64)
		if err != nil {
			return nil, err
		}

		ne.Number = num
	}

	return ne.Number, nil
}

func (ne NumberExpr) String() string {
	return fmt.Sprint(ne.NumStr)
}
