package lexer

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/spy16/parens/lexer/utfstrings"
)

var numberRegex = regexp.MustCompile("^(\\+|-)?\\d+(\\.\\d+)?$")

// scanComment advances the cursor till line-break or eof.
func scanComment(cur *utfstrings.Cursor) {
	for {
		ru := cur.Next()
		if ru == '\n' || ru == '\r' || ru == utfstrings.EOS {
			break
		}
	}
}

// scanNumber advances the cursor till a delimiting character
// is reached and collects all characters. If the collected
// characters don't match a number regex, returns false.
func scanNumber(cur *utfstrings.Cursor) bool {
	numStr := ""
	for {
		ru := cur.Next()
		if isSepratingChar(ru) || ru == utfstrings.EOS {
			cur.Backup()
			break
		}
		numStr = fmt.Sprintf("%s%c", numStr, ru)
	}
	if numberRegex.MatchString(numStr) {
		return true
	}

	return false
}

// scanSymbol advances the cursor until delimiting character is
// reached or a invalid rune is reached. resets the cursor position
// in case of invalid rune.
func scanSymbol(cur *utfstrings.Cursor) bool {
	runes := []rune{}
	for {
		ru := cur.Next()
		if ru == utfstrings.EOS {
			break
		} else if isSepratingChar(ru) {
			cur.Backup()
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

func scanString(cur *utfstrings.Cursor) error {
	cur.Next() // consume double-quote

	for {
		ru := cur.Next()
		if ru == '\\' {
			nextRune := cur.Peek()
			if nextRune == '"' || nextRune == 't' || nextRune == 'n' || nextRune == 'r' {
				cur.Next()
			}
		}

		if ru == '"' {
			return nil
		}

		if ru == utfstrings.EOS {
			return ErrUnterminatedString
		}

	}
}

// scanInvalidToken scans the current unidentified token and returns
// an error.
func scanInvalidToken(cur *utfstrings.Cursor) error {
	oldSel := cur.Selection
	unrec := ""
	for {
		ru := cur.Next()
		if ru == utfstrings.EOS || isSepratingChar(ru) {
			cur.Selection = oldSel
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
