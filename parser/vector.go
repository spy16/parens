package parser

// VectorExpr represents a vector form.
type VectorExpr struct {
	vector []Expr
}

// Eval creates a golang slice.
func (ve *VectorExpr) Eval(scope Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, expr := range ve.vector {
		val, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}
