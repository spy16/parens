package parser

import "github.com/spy16/parens/lexer"

type tokenQueue struct {
	tokens []lexer.Token
}

func (tq *tokenQueue) Token(index int) *lexer.Token {
	if index >= len(tq.tokens) {
		return nil
	}

	return &(tq.tokens[index])
}

// Pop removes and returns the the top element (i.e., 0th index) in the queue.
func (tq *tokenQueue) Pop() *lexer.Token {
	if len(tq.tokens) == 0 {
		return nil
	}

	token := tq.tokens[0]
	tq.tokens = tq.tokens[1:]

	if token.Type == lexer.WHITESPACE || token.Type == lexer.NEWLINE ||
		token.Type == lexer.COMMENT {
		return tq.Pop()
	}
	return &token
}
