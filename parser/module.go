package parser

import (
	"fmt"
	"strings"
)

// ModuleExpr represents a list of Exprs.
type ModuleExpr struct {
	Name  string
	Exprs []Expr
}

// Eval executes each expression in the module and returns the last result.
func (me ModuleExpr) Eval(scope Scope) (interface{}, error) {
	var val interface{}
	var err error

	for _, expr := range me.Exprs {
		val, err = expr.Eval(scope)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

func (me ModuleExpr) String() string {
	strs := []string{}
	for _, expr := range me.Exprs {
		strs = append(strs, fmt.Sprint(expr))
	}
	return strings.Join(strs, "\n")
}

func buildModuleExpr(name string, queue *tokenQueue) (Expr, error) {
	me := ModuleExpr{}
	me.Name = name

	for {
		expr, err := buildExpr(queue)
		if err != nil {
			if err == ErrEOF {
				break
			}
			return nil, err
		}

		if expr != nil {
			me.Exprs = append(me.Exprs, expr)
		}

	}

	return me, nil
}
