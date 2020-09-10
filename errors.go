package parens

import (
	"fmt"
	"reflect"

	"github.com/spy16/parens/value"
)

// Error is a root interface shared by all top-level parens errors.
type Error interface {
	ParensError()
}

// SyntaxError indicates an invalid form or expression.
type SyntaxError struct{ sentinel }

// NewSyntaxError wraps the supplied error.
func NewSyntaxError(err error) SyntaxError {
	return SyntaxError{sentinel{err}}
}

func (err SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s", err.error)
}

// Unwrap provides compatibility with go1.13 error chains
func (err SyntaxError) Unwrap() error { return err.error }

// TypeError is returned by the evaluator if an expression received an incorrect type.
type TypeError struct {
	Type reflect.Type
	sentinel
}

// NewTypeError wraps the error in a TypeError, annotating it with the offending value's
// type information.
func NewTypeError(offender value.Any, err error) TypeError {
	return TypeError{
		Type:     reflect.TypeOf(offender),
		sentinel: sentinel{err},
	}
}

// // Prefix in Error() message
// func (TypeError) Prefix() string { return "TypeError" }

func (err TypeError) Error() string {
	return fmt.Sprintf("TypeError: '%s' %s", err.Type, err.error)
}

// Unwrap provides compatibility with go1.13 error chains
func (err TypeError) Unwrap() error { return err.error }

// sentinel is a convenience struct that can be embedded in public errors to provide
// them with the ParensError() sentinel method (see:  repl/printer.go).
type sentinel struct{ error }

func (sentinel) ParensError() {}
