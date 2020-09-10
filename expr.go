package parens

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spy16/parens/value"
)

var (
	_ Expr = (*ConstExpr)(nil)
	_ Expr = (*DefExpr)(nil)
	_ Expr = (*QuoteExpr)(nil)
	_ Expr = (*InvokeExpr)(nil)
	_ Expr = (*IfExpr)(nil)
	_ Expr = (*DoExpr)(nil)
)

var (
	// ErrInvalidBindName is returned by DefExpr when the bind name is invalid.
	ErrInvalidBindName = errors.New("invalid name for def")

	// ErrNotInvokable is returned by InvokeExpr when the target is not invokable.
	ErrNotInvokable = errors.New("not invokable")
)

// ConstExpr returns the Const value wrapped inside when evaluated. It has
// no side-effect on the VM.
type ConstExpr struct{ Const value.Any }

// Eval the expression
func (ce ConstExpr) Eval(Context, Evaluator) (value.Any, error) { return ce.Const, nil }

// QuoteExpr expression represents a quoted form and
type QuoteExpr struct{ Form value.Any }

// Eval the expression
func (qe QuoteExpr) Eval(Context, Evaluator) (value.Any, error) {
	// TODO: re-use this for syntax-quote and unquote?
	return qe.Form, nil
}

// DefExpr creates a global binding with the Name when evaluated.
type DefExpr struct {
	Name  string
	Value value.Any
}

// Eval the expression
func (de DefExpr) Eval(ctx Context, ev Evaluator) (value.Any, error) {
	de.Name = strings.TrimSpace(de.Name)
	if de.Name == "" {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidBindName, de.Name)
	}

	ctx.SetGlobal(de.Name, de.Value)
	return &value.Symbol{Value: de.Name}, nil
}

// IfExpr represents the if-then-else form.
type IfExpr struct{ Test, Then, Else value.Any }

// Eval the expression
func (ife IfExpr) Eval(ctx Context, ev Evaluator) (value.Any, error) {
	test, err := ev.Eval(ctx, ife.Test)
	if err != nil {
		return nil, err
	}
	if value.IsTruthy(test) {
		return ev.Eval(ctx, ife.Then)
	}
	return ev.Eval(ctx, ife.Else)
}

// DoExpr represents the (do expr*) form.
type DoExpr struct{ Forms []value.Any }

// Eval the expression
func (de DoExpr) Eval(ctx Context, ev Evaluator) (value.Any, error) {
	var res value.Any
	var err error

	for _, form := range de.Forms {
		res, err = ev.Eval(ctx, form)
		if err != nil {
			return nil, err
		}
	}

	if res == nil {
		return value.Nil{}, nil
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
func (ie InvokeExpr) Eval(ctx Context, ev Evaluator) (value.Any, error) {
	val, err := ie.Target.Eval(ctx, ev)
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Invokable)
	if !ok {
		return nil, NewTypeError(val, ErrNotInvokable)
	}

	var args []value.Any
	for _, ae := range ie.Args {
		v, err := ae.Eval(ctx, ev)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}

	ctx.Push(StackFrame{Name: ie.Name, Args: args})
	defer ctx.Pop()

	return fn.Invoke(ev, args...)
}

// GoExpr evaluates a set of forms in a separate goroutine.
type GoExpr struct {
	Fn value.Seq
}

// Eval the expression.
func (ge GoExpr) Eval(ctx Context, ev Evaluator) (value.Any, error) {
	go ev.Eval(ctx.NewChild(), ge.Fn)
	return nil, nil
}
