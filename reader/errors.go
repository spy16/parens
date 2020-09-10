package reader

import (
	"errors"
	"fmt"
)

var (
	// ErrSkip can be returned by reader macro to indicate a no-op form which
	// should be discarded (e.g., Comments).
	ErrSkip = errors.New("skip expr")

	// ErrEOF is returned by reader when stream ends prematurely to indicate
	// that more data is needed to complete the current form.
	ErrEOF = errors.New("unexpected EOF")
)

// UnmatchedDelimiterError is returned when a reader macro encounters a closing container-
// delimiter without a corresponding opening delimiter (e.g. ']' but no '[').
type UnmatchedDelimiterError rune

func (err UnmatchedDelimiterError) Error() string {
	return fmt.Sprintf("unmatched delimiter: '%c'", err)
}

// NumberFormatError is returned when a reader macro encounters a illegally-formatted
// numerical form.
type NumberFormatError string

func (err NumberFormatError) Error() string {
	return fmt.Sprintf("illegal number format: '%s'", string(err))
}

// ReadError is a generic error that is returned when a reader macro fails to read from
// the character stream.
type ReadError struct {
	FormType string
	Value    error
}

// NewReadErrorFactory binds a formType to a ReadError constructor.
// It is provided as a conveience for generating errors for a specific type of form.
func NewReadErrorFactory(formType string) func(error) ReadError {
	return func(err error) ReadError {
		return ReadError{
			FormType: formType,
			Value:    err,
		}
	}
}

func (err ReadError) Error() string {
	if err.Value == nil {
		return fmt.Sprintf("error reading %s form", err.FormType)
	}

	return fmt.Sprintf("error reading %s form: %s", err.FormType, err.Value)
}

// Unwrap provides compatibility with go1.13 error chains.
func (err ReadError) Unwrap() error {
	return err.Value
}

// SyntaxError .
type SyntaxError struct {
	Cause error
}

func (err SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s", err.Cause)
}

// Unwrap provides compatibility with go1.13 error chains
func (err SyntaxError) Unwrap() error {
	return err.Cause
}
