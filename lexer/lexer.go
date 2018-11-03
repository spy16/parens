package lexer

import (
	"errors"
	"unicode/utf8"
)

const eof = -1

// ErrEOF is returned when the lexer reaches the end of file.
var ErrEOF = errors.New("end of file")

// New initializes the lexer with the given source. Source can
// contain any UTF-8 characters.
func New(src string) *Lexer {
	return &Lexer{
		src: src,
	}
}

// Lexer performs lexical analysis of LISP.
type Lexer struct {
	cursor

	src string
}

type cursor struct {
	pos   int
	start int
	width int
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
	token.Start = lex.start
	token.Value = lex.src[lex.start:lex.pos]
	token.Type = tokenType

	lex.start = lex.pos
	return &token, nil
}

func (lex *Lexer) nextTokenType() (TokenType, error) {
	switch ru := lex.next(); {
	case ru == eof:
		return "", ErrEOF

	case ru == '\'':
		return QUOTE, nil

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
		lex.backup()
		if err := scanString(lex); err != nil {
			return "", err
		}
		return STRING, nil

	case ru == ';':
		lex.backup()
		scanComment(lex)
		return COMMENT, nil

	default:
		lex.backup()
		oldCursor := lex.cursor
		if scanNumber(lex) {
			return NUMBER, nil
		}
		lex.cursor = oldCursor

		if scanSymbol(lex) {
			return SYMBOL, nil
		}
		lex.cursor = oldCursor

		return "", scanInvalidToken(lex)
	}
}

// next returns the next rune in the input.
func (lex *Lexer) next() rune {
	if int(lex.pos) >= len(lex.src) {
		lex.width = 0
		return eof
	}
	ru, width := utf8.DecodeRuneInString(lex.src[lex.pos:])
	lex.width = width
	lex.pos += lex.width
	return ru
}

// peek returns but does not consume the next rune in the input.
func (lex *Lexer) peek() rune {
	r := lex.next()
	lex.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (lex *Lexer) backup() {
	lex.pos -= lex.width
}
