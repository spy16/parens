package parens

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a binding not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidBindName is returned by DefExpr when the bind name is invalid.
	ErrInvalidBindName = errors.New("invalid name for def")

	// ErrNotInvokable is returned by InvokeExpr when the target is not invokable.
	ErrNotInvokable = errors.New("not invokable")

	// ErrIncomparableTypes is returned by Any.Comp when a comparison between two tpyes
	// is undefined.  Users should generally consider the types to be not equal in such
	// cases, but not assume any ordering.
	ErrIncomparableTypes = errors.New("incomparable types")
)

// New returns a new root context initialised based on given options.
func New(opts ...Option) *Env {
	env := &Env{ctx: context.Background(), globals: newMutexMap()}
	for _, opt := range withDefaults(opts) {
		opt(env)
	}
	return env
}

// Any represents any
type Any interface {
	// SExpr MUST return a parsable s-expression that can be consumed by
	// a reader.Reader.
	//
	// For a human-readable implementation, implement `repl.Renderable`.
	SExpr() (string, error)
}

// Seq represents a sequence of values.
type Seq interface {
	Any
	Count() (int, error)
	First() (Any, error)
	Next() (Seq, error)
	Conj(items ...Any) (Seq, error)
}

// Analyzer implementation is responsible for performing syntax analysis
// on given form.
type Analyzer interface {
	// Analyze should perform syntax checks for special forms etc. and
	// return Expr values that can be evaluated against a context.
	Analyze(env *Env, form Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(env *Env, form Any) (Any, error)
}

// Invokable represents a value that can be invoked for result.
type Invokable interface {
	Invoke(env *Env, args ...Any) (Any, error)
}

// Expr represents an expression that can be evaluated against a context.
type Expr interface {
	Eval() (Any, error)
}

// Error is returned by all parens operations. Cause indicates the underlying
// error type. Use errors.Is() with Cause to check for specific errors.
type Error struct {
	Message string
	Cause   error
}

// Is returns true if the other error is same as the cause of this error.
func (e Error) Is(other error) bool { return errors.Is(e.Cause, other) }

// Unwrap returns the underlying cause of the error.
func (e Error) Unwrap() error { return e.Cause }

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %s", e.Cause, e.Message)
	}
	return e.Message
}
