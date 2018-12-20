package parens

import (
	"fmt"
	"io"
)

// Execute reads until EOF or an error from the RuneScanner and executes the
// read s-expressions in the given scope.
func Execute(name string, rd io.RuneScanner, env Scope) (interface{}, error) {
	expr, err := ParseModule(name, rd)
	if err != nil {
		return nil, err
	}

	return ExecuteExpr(expr, env)
}

// ExecuteOne reads runes enough to construct one s-exp and executes the s-exp
// with given scope.
func ExecuteOne(rd io.RuneScanner, env Scope) (interface{}, error) {
	expr, err := Parse(rd)
	if err != nil {
		return nil, err
	}

	return ExecuteExpr(expr, env)
}

// ExecuteExpr executes the expr in the given scope.
func ExecuteExpr(expr Expr, env Scope) (interface{}, error) {
	var res interface{}
	var evalErr error
	safeWrapper := func() {
		defer func() {
			if v := recover(); v != nil {
				if err, ok := v.(error); ok {
					evalErr = err
				} else {
					evalErr = fmt.Errorf("panic: %v", v)
				}
			}
		}()

		res, evalErr = expr.Eval(env)
	}

	safeWrapper()
	if evalErr != nil {
		return nil, evalErr
	}

	return res, nil
}
