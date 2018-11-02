package lexer

import (
	"fmt"
)

// Token represents an extracted token from the stream
type Token struct {
	Type  TokenType
	Start int
	Value string
}

func (token Token) String() string {
	return fmt.Sprintf("Token{type=%s, position=%d,%d}", token.Type, token.Start, token.Start+len(token.Value))
}

// TokenType represents the type of the extracted token
type TokenType string

const (
	//UNKNOWN represents a state where the token is yet to
	// be identified
	UNKNOWN TokenType = ""

	// LPAREN represents the left-parenthesis character
	LPAREN TokenType = "LPAREN"

	// RPAREN represents the right-parenthesis character
	RPAREN TokenType = "RPAREN"

	// SSTRING represents a single-quoted string
	SSTRING TokenType = "SSTRING"

	// DSTRING represents a double-quoted string
	DSTRING TokenType = "DSTRING"

	// NUMBER represents a int/float number
	NUMBER TokenType = "NUMBER"

	// WHITESPACE represents a space, tab or newline
	WHITESPACE TokenType = "WHITESPACE"

	// COMMENT represents a ";" based comment
	COMMENT TokenType = "COMMENT"

	// SYMBOL represents any identifier
	SYMBOL TokenType = "SYMBOL"
)
