package lexer

import (
	"errors"

	"github.com/spy16/parens/lexer/utfstrings"
)

// ErrEOF is returned when the lexer reaches the end of file.
var ErrEOF = errors.New("end of file")

// New initializes the lexer with the given source. Source can
// contain any UTF-8 characters.
func New(src string) *Lexer {
	return &Lexer{
		cur: utfstrings.Cursor{
			String: src,
		},
	}
}

// Lexer performs lexical analysis of LISP.
type Lexer struct {
	cur utfstrings.Cursor
}

// Tokens runs through the entire source and returns tokens.
func (lex *Lexer) Tokens() ([]Token, error) {
	tokens := []Token{}
	for {
		token, err := lex.Next()
		if err != nil {
			if err == ErrEOF {
				return tokens, nil
			}
			return nil, err
		}

		tokens = append(tokens, *token)
	}
}

// Next consumes characters from source until the next token is found
// and returns the token. If no token is identified till the end of
// source, ErrEOF will be returned.
func (lex *Lexer) Next() (*Token, error) {
	tokenType, err := lex.nextTokenType()
	if err != nil {
		return nil, err
	}

	token := Token{}
	token.Start = lex.cur.Start
	token.Value = lex.cur.String[lex.cur.Start:lex.cur.Pos]
	token.Type = tokenType

	lex.cur.Start = lex.cur.Pos
	return &token, nil
}

func (lex *Lexer) nextTokenType() (TokenType, error) {
	switch ru := lex.cur.Next(); {
	case ru == utfstrings.EOS:
		return "", ErrEOF

	case ru == '(':
		return LPAREN, nil

	case ru == ')':
		return RPAREN, nil

	case ru == '[':
		return LVECT, nil

	case ru == ']':
		return RVECT, nil

	case ru == '\n' || ru == '\r':
		return NEWLINE, nil

	case ru == ' ' || ru == '\t':
		return WHITESPACE, nil

	case ru == '"':
		lex.cur.Backup()
		if err := scanString(&lex.cur); err != nil {
			return "", err
		}
		return STRING, nil

	case ru == ';':
		lex.cur.Backup()
		scanComment(&lex.cur)
		return COMMENT, nil

	case ru == '\'':
		return QUOTE, nil

	default:
		lex.cur.Backup()
		oldSel := lex.cur.Selection
		if scanNumber(&lex.cur) {
			return NUMBER, nil
		}
		lex.cur.Selection = oldSel

		if scanSymbol(&lex.cur) {
			return SYMBOL, nil
		}
		lex.cur.Selection = oldSel

		return "", scanInvalidToken(&lex.cur)
	}
}
