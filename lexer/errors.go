package lexer

import (
	"errors"
	"fmt"
)

// ErrUnterminatedString is returned when an open-quote does not have a matching
// close.
var ErrUnterminatedString = errors.New("unterminated string literal")

// ErrUnrecognizedToken is returned when a character or sequence
// of characters cannot be recognized as a valid token.
type ErrUnrecognizedToken struct {
	val string
}

func (err ErrUnrecognizedToken) Error() string {
	return fmt.Sprintf("unrecognized token '%s'", err.val)
}
