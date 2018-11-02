package lexer

import (
	"fmt"
	"unicode"
)

func isSepratingChar(ru rune) bool {
	return oneOf(ru, ' ', '\t', '\n', '\r', '(', ')', '[', ']')
}

func scanComment(lex *Lexer) {
	for {
		ru := lex.next()
		if ru == '\n' || ru == '\r' || ru == eof {
			break
		}
	}
}

func scanNumber(lex *Lexer) bool {
	oldCursor := lex.cursor

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

	lex.cursor = oldCursor
	return false
}

func scanSymbol(lex *Lexer) bool {
	oldCursor := lex.cursor

	for {
		switch ru := lex.next(); {
		case ru == eof:
			return true

		case isSepratingChar(ru):
			lex.backup()
			return true

		case !isAlphaNumeric(ru):
			lex.cursor = oldCursor
			return false
		}
	}
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

func oneOf(ru rune, set ...rune) bool {
	for _, rs := range set {
		if ru == rs {
			return true
		}
	}
	return false
}

// isAlphaNumeric reports whether r is a valid rune for an identifier.
func isAlphaNumeric(r rune) bool {
	return r == '>' || r == '<' || r == '=' || r == '-' || r == '+' || r == '*' || r == '&' || r == '_' || r == '/' || r == '?' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
