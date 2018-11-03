package parser

import (
	"errors"
	"fmt"

	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/reflection"
)

// ErrEOF is returned when the parser has consumed all tokens.
var ErrEOF = errors.New("end of file")

// Parse the slice of tokens and build an AST
func Parse(src string) (SExp, error) {
	tokens, err := lexer.New(src).Tokens()
	if err != nil {
		return nil, err
	}

	return buildSExp(&tokenQueue{tokens: tokens})
}

// SExp represents a symbolic expression.
type SExp interface {
	Eval(env *reflection.Env) (interface{}, error)
}

func buildSExp(tokens *tokenQueue) (SExp, error) {
	if len(tokens.tokens) == 0 {
		return nil, ErrEOF
	}

	token := tokens.Pop()

	switch token.Type {
	case lexer.LPAREN:
		le := ListExp{}

		for {
			next := tokens.Token(0)
			if next == nil {
				return nil, ErrEOF
			}
			if next.Type == lexer.RPAREN {
				break
			}
			exp, err := buildSExp(tokens)
			if err != nil {
				return nil, err
			}

			if exp == nil {
				continue
			}

			le.List = append(le.List, exp)
		}
		tokens.Pop()
		return le, nil

	case lexer.RPAREN:
		return nil, ErrEOF

	case lexer.NUMBER:
		ne := NumberExp{
			Token: *token,
		}
		return ne, nil

	case lexer.STRING:
		se := StringExp{
			Token: *token,
		}
		return se, nil

	case lexer.SYMBOL:
		se := SymbolExp{
			Symbol: token.Value,
		}
		return se, nil

	case lexer.WHITESPACE, lexer.NEWLINE:
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown token type: %s", (token.Type))
	}
}
