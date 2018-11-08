package parser

import (
	"errors"
	"fmt"

	"github.com/spy16/parens/lexer"
)

// ErrEOF is returned when the parser has consumed all tokens.
var ErrEOF = errors.New("end of file")

// Parse tokenizes and parses the src to build an AST.
func Parse(name string, src string) (Expr, error) {
	tokens, err := lexer.New(src).Tokens()
	if err != nil {
		return nil, err
	}

	queue := &tokenQueue{tokens: tokens}
	return buildModuleExpr(name, queue)
}

// Expr represents an evaluatable expression.
type Expr interface {
	Eval(env Scope) (interface{}, error)
}

// Scope is responsible for managing bindings.
type Scope interface {
	Get(name string) (interface{}, error)
	Doc(name string) string
	Bind(name string, v interface{}, doc ...string) error
	Root() Scope
}

func buildExpr(tokens *tokenQueue) (Expr, error) {
	if len(tokens.tokens) == 0 {
		return nil, ErrEOF
	}

	token := tokens.Pop()

	switch token.Type {
	case lexer.LPAREN:
		return buildListExpr(tokens)

	case lexer.NUMBER:
		return newNumberExpr(token), nil

	case lexer.STRING:
		return newStringExpr(token), nil

	case lexer.SYMBOL:
		return newSymbolExpr(token), nil

	case lexer.LVECT:
		return buildVectorExpr(tokens)

	case lexer.QUOTE:
		expr, err := buildExpr(tokens)
		if err != nil {
			return nil, err
		}
		return QuoteExpr{expr: expr}, nil

	case lexer.WHITESPACE, lexer.NEWLINE, lexer.COMMENT:
		return nil, nil

	case lexer.RPAREN, lexer.RVECT:
		return nil, ErrEOF

	default:
		return nil, fmt.Errorf("unknown token type: %s", (token.Type))
	}

}
