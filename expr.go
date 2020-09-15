package parens

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	_ Expr = (*ConstExpr)(nil)
	_ Expr = (*DefExpr)(nil)
	_ Expr = (*QuoteExpr)(nil)
	_ Expr = (*InvokeExpr)(nil)
	_ Expr = (*IfExpr)(nil)
	_ Expr = (*DoExpr)(nil)
)

// ConstExpr returns the Const value wrapped inside when evaluated. It has
// no side-effect on the VM.
type ConstExpr struct{ Const Any }

// Eval returns the constant value unmodified.
func (ce ConstExpr) Eval(_ *Env) (Any, error) { return ce.Const, nil }

// QuoteExpr expression represents a quoted form and
type QuoteExpr struct{ Form Any }

// Eval returns the quoted form unmodified.
func (qe QuoteExpr) Eval(_ *Env) (Any, error) {
	// TODO: re-use this for syntax-quote and unquote?
	return qe.Form, nil
}

// DefExpr creates a global binding with the Name when evaluated.
type DefExpr struct {
	Name  string
	Value Any
}

// Eval creates a symbol binding in the global (root) stack frame.
func (de DefExpr) Eval(env *Env) (Any, error) {
	de.Name = strings.TrimSpace(de.Name)
	if de.Name == "" {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidBindName, de.Name)
	}

	env.setGlobal(de.Name, de.Value)
	return Symbol(de.Name), nil
}

// IfExpr represents the if-then-else form.
type IfExpr struct{ Test, Then, Else Any }

// Eval the expression
func (ife IfExpr) Eval(env *Env) (Any, error) {
	test, err := env.Eval(ife.Test)
	if err != nil {
		return nil, err
	}
	if IsTruthy(test) {
		return env.Eval(ife.Then)
	}
	return env.Eval(ife.Else)
}

// DoExpr represents the (do expr*) form.
type DoExpr struct{ Forms []Any }

// Eval the expression
func (de DoExpr) Eval(env *Env) (Any, error) {
	var res Any
	var err error

	for _, form := range de.Forms {
		res, err = env.Eval(form)
		if err != nil {
			return nil, err
		}
	}

	if res == nil {
		return Nil{}, nil
	}
	return res, nil
}

// InvokeExpr performs invocation of target when evaluated.
type InvokeExpr struct {
	Name   string
	Target Expr
	Args   []Expr
}

// Eval the expression
func (ie InvokeExpr) Eval(env *Env) (Any, error) {
	val, err := ie.Target.Eval(env)
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Invokable)
	if !ok {
		return nil, Error{
			Cause:   ErrNotInvokable,
			Message: fmt.Sprintf("value of type '%s' is not invokable", reflect.TypeOf(val)),
		}
	}

	var args []Any
	for _, ae := range ie.Args {
		v, err := ae.Eval(env)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}

	env.push(stackFrame{
		Name: ie.Name,
		Args: args,
		Vars: map[string]Any{},
	})
	defer env.pop()

	return fn.Invoke(env, args...)
}

// GoExpr evaluates an expression in a separate goroutine.
type GoExpr struct {
	Value Any
}

// Eval forks the given context to get a child context and launches goroutine
// with the child context to evaluate the
func (ge GoExpr) Eval(env *Env) (Any, error) {
	child := env.fork()
	go func() {
		_, _ = child.Eval(ge.Value)
	}()
	return nil, nil
}
