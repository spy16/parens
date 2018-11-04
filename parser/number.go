package parser

import (
	"strconv"

	"github.com/spy16/parens/reflection"
)

// NumberExp represents number s-expression.
type NumberExp struct {
	numStr string
	number *float64
}

// Eval for a number returns itself.
func (ne NumberExp) Eval(scope *reflection.Scope) (interface{}, error) {
	if ne.number == nil {
		num, err := strconv.ParseFloat(ne.numStr, 64)
		if err != nil {
			return nil, err
		}

		ne.number = &num
	}

	return *ne.number, nil
}
