package parser

import (
	"github.com/spy16/parens/reflection"
)

// VectorExp represents a vector form.
type VectorExp struct {
	Vector []SExp
}

// Eval creates a golang slice.
func (ve *VectorExp) Eval(scope *reflection.Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, sexp := range ve.Vector {
		val, err := sexp.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}
