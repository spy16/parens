package parens

import (
	"strconv"
)

// MacroFunc represents the signature of the Go macro functions. Functions
// bound in the scope as MacroFunc will receive un-evaluated list of s-exps
// and the current scope.
type MacroFunc func(scope Scope, exprs []Expr) (interface{}, error)

// ModuleExpr represents a list of Exprs.
type ModuleExpr []Expr

// Eval executes each expression in the module and returns the last result.
func (me ModuleExpr) Eval(scope Scope) (interface{}, error) {
	var val interface{}
	var err error

	for _, expr := range me {
		val, err = expr.Eval(scope)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

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

// StringExpr represents single and double quoted strings.
type StringExpr string

// Eval returns unquoted version of the STRING token.
func (se StringExpr) Eval(_ Scope) (interface{}, error) {
	return string(se), nil
}

// QuoteExpr implements the quote-literal form.
type QuoteExpr struct {
	Expr Expr
}

// Eval returns the expression itself without evaluating it.
func (qe QuoteExpr) Eval(scope Scope) (interface{}, error) {
	return qe.Expr, nil
}

// UnquoteEval unquotes and evaluates the underlying expression.
func (qe QuoteExpr) UnquoteEval(scope Scope) (interface{}, error) {
	return qe.Expr.Eval(scope)
}

// CommentExpr is returned to represent a lisp-style comment.
type CommentExpr string

// Eval returns the comment string.
func (ce CommentExpr) Eval(_ Scope) (interface{}, error) {
	return string(ce), nil
}

// KeywordExpr represents a keyword literal.
type KeywordExpr string

// Eval returns the keyword itself.
func (ke KeywordExpr) Eval(_ Scope) (interface{}, error) {
	return string(ke), nil
}

// SymbolExpr represents a symbol.
type SymbolExpr string

// Eval returns the symbol name itself.
func (se SymbolExpr) Eval(scope Scope) (interface{}, error) {
	return scope.Get(string(se))
}

// VectorExpr represents a vector form.
type VectorExpr []Expr

// Eval creates a golang slice.
func (ve VectorExpr) Eval(scope Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, expr := range ve {
		val, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}

// ListExpr represents a list (i.e., a function call) expression.
type ListExpr []Expr

// Eval evaluates each s-exp in the list and then evaluates the list itself
// as an s-exp.
func (le ListExpr) Eval(scope Scope) (interface{}, error) {
	if len(le) == 0 {
		return le, nil
	}

	val, err := le[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if macroFn, ok := val.(MacroFunc); ok {
		return macroFn(scope, le[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(le); i++ {
		arg, err := le[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return reflectCall(val, args...)
}
