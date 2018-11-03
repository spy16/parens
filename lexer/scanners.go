package lexer

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

var numberRegex = regexp.MustCompile("^(\\+|-)?\\d+(\\.\\d+)?$")

// scanComment advances the lexer till line-break or eof.
func scanComment(lex *Lexer) {
	for {
		ru := lex.next()
		if ru == '\n' || ru == '\r' || ru == eof {
			break
		}
	}
}

// scanNumber advances the lexer till a delimiting character
// is reached and collects all characters. If the collected
// characters don't match a number regex, returns false.
func scanNumber(lex *Lexer) bool {
	numStr := ""
	for {
		ru := lex.next()
		if isSepratingChar(ru) || ru == eof {
			lex.backup()
			break
		}
		numStr = fmt.Sprintf("%s%c", numStr, ru)
	}
	if numberRegex.MatchString(numStr) {
		return true
	}

	return false
}

// scanSymbol advances the lexer until delimiting character is
// reached or a invalid rune is reached. resets the lexer position
// in case of invalid rune.
func scanSymbol(lex *Lexer) bool {
	runes := []rune{}
	for {
		ru := lex.next()
		if ru == eof {
			break
		} else if isSepratingChar(ru) {
			lex.backup()
			break
		} else if !utf8.ValidRune(ru) {
			return false
		}
		runes = append(runes, ru)
	}

	if unicode.IsDigit(runes[0]) {
		return false
	}

	return true
}

func scanString(lex *Lexer) error {
	lex.next() // consume double-quote

	for {
		ru := lex.next()
		if ru == '\\' {
			nextRune := lex.peek()
			if nextRune == '"' || nextRune == 't' || nextRune == 'n' || nextRune == 'r' {
				lex.next()
			}
		}

		if ru == '"' {
			return nil
		}

		if ru == eof {
			return fmt.Errorf("unterminated string")
		}

	}
}

// scanInvalidToken scans the current unidentified token and returns
// an error.
func scanInvalidToken(lex *Lexer) error {
	oldCursor := lex.cursor
	unrec := ""
	for {
		ru := lex.next()
		if ru == eof || isSepratingChar(ru) {
			lex.cursor = oldCursor
			break
		}

		unrec = fmt.Sprintf("%s%c", unrec, ru)
	}

	return &ErrUnrecognizedToken{val: unrec}
}

func oneOf(ru rune, set ...rune) bool {
	for _, rs := range set {
		if ru == rs {
			return true
		}
	}
	return false
}

func isSepratingChar(ru rune) bool {
	return oneOf(ru, ' ', '\t', '\n', '\r', '(', ')', '[', ']')
}
