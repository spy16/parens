package lexer

import (
	"errors"
	"fmt"
	"regexp"
)

var errNoToken = errors.New("no token identified")

// consumer should consume bytes from src until a token is found.
// If token is not found, consumer should return errNoToken. if
// a token is malformed, consumer should return appropriate errors.
// if token is found, consumer should return a Token type with
// `Type` and `Value` initialized.
type consumer func(src []byte) (*Token, error)

func newNumberConsumer() consumer {
	return func(src []byte) (*Token, error) {
		t := &Token{
			Type:  NUMBER,
			Value: "",
		}

		decimalPointFound := false
		atleastOneDigitFound := false
		i := 0
		for i < len(src) {
			crune := src[i]
			if crune == '.' {
				if decimalPointFound {
					break
				}
				decimalPointFound = true
				t.Value += "."
			} else if crune == '-' {
				if i != 0 {
					break
				}
				t.Value = "-"
			} else if crune >= 48 && crune <= 57 {
				t.Value += string(crune)
				atleastOneDigitFound = true
			} else {
				break
			}

			i++
		}

		if t.Value == "" || atleastOneDigitFound == false {
			return nil, nil
		}

		return t, nil
	}
}

func newQuotedStringConsumer(typ TokenType, quote byte) consumer {
	return func(src []byte) (*Token, error) {
		var token Token
		token.Type = typ

		if len(src) <= 0 {
			// no characters left to consume
			return nil, nil
		}

		if src[0] != quote {
			// not a quoted string
			return nil, nil
		}

		if len(src) < 2 {
			return nil, fmt.Errorf("unexpected EOF, expecting '%s'", string(quote))
		}

		i := 1
		inEscape := false
		for i < len(src) {
			if src[i] == quote && !inEscape {
				break
			}
			inEscape = false

			if src[i] == '\\' {
				inEscape = true
			}

			if i == len(src)-1 {
				return nil, fmt.Errorf("unexpected EOF, expecting '%s'", string(quote))
			}

			i++
		}

		token.Value = string(src[:i+1])

		return &token, nil

	}
}

func newRegExpConsumer(typ TokenType, exp string) consumer {
	rxp := regexp.MustCompile(exp)
	return func(src []byte) (*Token, error) {
		matches := rxp.FindSubmatch(src)
		if len(matches) > 1 {
			token := Token{
				Type:  typ,
				Value: string(matches[1]),
			}

			return &token, nil
		}
		return nil, nil
	}
}
