package lexer

import (
	"fmt"
)

// ErrUnrecognizedToken is returned when a character or sequence
// of characters cannot be recognized as a valid token.
type ErrUnrecognizedToken struct {
	val rune
}

func (err ErrUnrecognizedToken) Error() string {
	return fmt.Sprintf("unrecognized token '%s'", string(err.val))
}
