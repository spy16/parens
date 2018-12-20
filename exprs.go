package parens

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MacroFunc represents the signature of the Go macro functions. Functions
// bound in the scope as MacroFunc will receive un-evaluated list of s-exps
// and the current scope.
type MacroFunc func(scope Scope, exprs []Expr) (interface{}, error)

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

// StringExpr represents single and double quoted strings.
type StringExpr struct {
	Value string
}

// Eval returns unquoted version of the STRING token.
func (se StringExpr) Eval(_ Scope) (interface{}, error) {
	return se.Value, nil
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

func (qe QuoteExpr) String() string {
	return fmt.Sprintf("'%s", qe.Expr)
}

// CommentExpr is returned to represent a lisp-style comment.
type CommentExpr struct {
	comment string
}

// Eval returns the comment string.
func (ce CommentExpr) Eval(_ Scope) (interface{}, error) {
	return ce.comment, nil
}

func (ce CommentExpr) String() string {
	return ce.comment
}

// KeywordExpr represents a keyword literal.
type KeywordExpr struct {
	Keyword string
}

// Eval returns the keyword itself.
func (ke KeywordExpr) Eval(_ Scope) (interface{}, error) {
	return ke.Keyword, nil
}

func (ke KeywordExpr) String() string {
	return ke.Keyword
}

// SymbolExpr represents a symbol.
type SymbolExpr struct {
	Symbol string
}

// Eval returns the symbol name itself.
func (se SymbolExpr) Eval(scope Scope) (interface{}, error) {
	parts := strings.Split(se.Symbol, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid member access symbol. must be of format <parent>.<member>")
	}

	obj, err := scope.Get(parts[0])
	if err != nil {
		return nil, err
	}

	if len(parts) == 1 {
		return obj, nil
	}

	member := resolveMember(reflect.ValueOf(obj), parts[1])
	if !member.IsValid() {
		return nil, fmt.Errorf("member '%s' not found on '%s'", parts[1], parts[0])
	}

	return member.Interface(), nil
}

func (se SymbolExpr) String() string {
	return se.Symbol
}

func resolveMember(obj reflect.Value, name string) reflect.Value {
	firstMatch := func(fxs ...func(string) reflect.Value) reflect.Value {
		for _, fx := range fxs {
			if val := fx(name); val.IsValid() && val.CanInterface() {
				return val
			}
		}

		return reflect.Value{}
	}

	var funcs []func(string) reflect.Value
	if obj.Kind() == reflect.Ptr {
		funcs = append(funcs,
			obj.Elem().FieldByName,
			obj.MethodByName,
			obj.Elem().MethodByName,
		)
	} else {
		funcs = append(funcs,
			obj.FieldByName,
			obj.MethodByName,
		)
	}

	return firstMatch(funcs...)
}

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

// ListExpr represents a list (i.e., a function call) expression.
type ListExpr struct {
	List []Expr
}

// Eval evaluates each s-exp in the list and then evaluates the list itself
// as an s-exp.
func (le ListExpr) Eval(scope Scope) (interface{}, error) {
	if len(le.List) == 0 {
		return le.List, nil
	}

	val, err := le.List[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if macroFn, ok := val.(MacroFunc); ok {
		return macroFn(scope, le.List[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(le.List); i++ {
		arg, err := le.List[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return Call(val, args...)
}

func (le ListExpr) String() string {
	reprs := []string{}
	for _, item := range le.List {
		reprs = append(reprs, fmt.Sprint(item))
	}

	return fmt.Sprintf("(%s)", strings.Join(reprs, " "))
}
