package lexer

import (
	"fmt"
)

// ErrUnrecognizedToken is returned when a character or sequence
// of characters cannot be recognized as a valid token.
type ErrUnrecognizedToken struct {
	val string
}

func (err ErrUnrecognizedToken) Error() string {
	return fmt.Sprintf("unrecognized token '%s'", err.val)
}
