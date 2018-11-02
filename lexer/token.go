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
	return fmt.Sprintf("Token{type=%s, position=%d:%d, value='%s'}",
		token.Type, token.Start, token.Start+len(token.Value), token.Value)
}

// TokenType represents the type of the extracted token
type TokenType string

const (
	// LPAREN represents the left-parenthesis character
	LPAREN TokenType = "LPAREN"

	// RPAREN represents the right-parenthesis character
	RPAREN TokenType = "RPAREN"

	// LVECT represents the left-brace character
	LVECT TokenType = "LVECT"

	// RVECT represents the right-brace character
	RVECT TokenType = "RVECT"

	// STRING represents a double-quoted string
	STRING TokenType = "STRING"

	// NUMBER represents int, float, hex, complex etc.
	NUMBER TokenType = "NUMBER"

	// WHITESPACE represents a space, tab or newline
	WHITESPACE TokenType = "WHITESPACE"

	// NEWLINE represents a new-line or return-line-feed character.
	NEWLINE TokenType = "NEWLINE"

	// COMMENT represents a ";" based comment
	COMMENT TokenType = "COMMENT"

	// SYMBOL represents any identifier
	SYMBOL TokenType = "SYMBOL"

	// QUOTE represents a single-quote
	QUOTE TokenType = "QUOTE"
)
