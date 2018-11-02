package lexer

import (
	"errors"
)

// ErrEOF represents the end-of-file
var ErrEOF = errors.New("EOF")

var consumers = []consumer{
	newRegExpConsumer(LPAREN, `^(\()`),
	newRegExpConsumer(RPAREN, `^(\))`),
	newQuotedStringConsumer(DSTRING, '"'),
	newNumberConsumer(),
	newRegExpConsumer(WHITESPACE, `^(\s+)`),
	newRegExpConsumer(SYMBOL, `^('|[^\s();]+)`),
}

// New initializes a new instance of lexer for the given
// source
func New(src []byte) *Lexer {
	lxr := &Lexer{}
	lxr.src = src
	lxr.pos = 0
	lxr.consumers = consumers
	return lxr
}

// Lexer represents an instance of lexer
type Lexer struct {
	src       []byte
	pos       int
	consumers []consumer
}

// Tokens returns a slice of all tokens
func (lxr *Lexer) Tokens() ([]Token, error) {
	ts := []Token{}
	for {
		t, err := lxr.NextToken()
		if err != nil {
			if err == ErrEOF {
				return ts, nil
			}
			return nil, err
		}
		ts = append(ts, *t)
	}
}

// NextToken consumes characters starting from current position until a
// valid token is identified. Finally returns the token.
func (lxr *Lexer) NextToken() (*Token, error) {
	if lxr.pos >= len(lxr.src) {
		return nil, ErrEOF
	}

	for _, consumerFunc := range consumers {
		token, err := consumerFunc(lxr.src[lxr.pos:])
		if err != nil {
			return nil, err
		}

		if token == nil {
			continue
		}

		token.Start = lxr.pos
		lxr.pos = token.Start + len(token.Value)
		return token, nil
	}
	return nil, newInvalidTokenError(lxr.pos)
}
