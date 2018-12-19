package parser

import (
	"fmt"
	"strings"
)

// VectorExpr represents a vector form.
type VectorExpr struct {
	List []Expr
}

// Eval creates a golang slice.
func (ve VectorExpr) Eval(scope Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, expr := range ve.List {
		val, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}

func (ve VectorExpr) String() string {
	strs := []string{}
	for _, expr := range ve.List {
		strs = append(strs, fmt.Sprint(expr))
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}
