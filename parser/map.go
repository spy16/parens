package parser

import (
	"fmt"

	"github.com/spy16/parens/lexer"
)

// MapExpr represents a map literal expression.
type MapExpr struct {
	hashMap map[string]Expr
}

// Eval evaluates a map literal expression into map[string]interface{}
func (me MapExpr) Eval(scope Scope) (interface{}, error) {
	m := map[string]interface{}{}
	for key, valExpr := range me.hashMap {
		val, err := valExpr.Eval(scope)
		if err != nil {
			return nil, err
		}

		m[key] = val
	}
	return m, nil
}

func buildMapExpr(queue *tokenQueue) (Expr, error) {
	me := MapExpr{}
	me.hashMap = map[string]Expr{}

	for {
		key := queue.Pop()
		if key == nil || key.Type == lexer.RDICT {
			break
		}

		if key.Type != lexer.KEYWORD {
			return nil, fmt.Errorf("expecting keyword, not '%s'", key.Type)
		}

		val, err := buildExpr(queue)
		if err != nil {
			return nil, err
		}
		me.hashMap[key.Value] = val
	}

	return me, nil
}
