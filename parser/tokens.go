package parser

import "github.com/spy16/parens/lexer"

type tokenQueue struct {
	tokens []lexer.Token
}

func (tq *tokenQueue) Token(index int) *lexer.Token {
	if index >= len(tq.tokens) {
		return nil
	}

	return &(tq.tokens[0])
}

func (tq *tokenQueue) TokenType(index int) *lexer.TokenType {
	if index >= len(tq.tokens) {
		return nil
	}

	return &(tq.tokens[0].Type)
}

// Pop removes and returns the the top element (i.e., 0th index) in the queue.
func (tq *tokenQueue) Pop() *lexer.Token {
	if len(tq.tokens) == 0 {
		return nil
	}

	token := tq.tokens[0]
	tq.tokens = tq.tokens[1:]
	return &token
}
