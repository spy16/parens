package lexer

// InvalidTokenError represents an invalid token
type InvalidTokenError struct {
	Pos int
}

func (ite InvalidTokenError) Error() string {
	return "invalid token"
}

func newInvalidTokenError(pos int) InvalidTokenError {
	return InvalidTokenError{
		Pos: pos,
	}
}
