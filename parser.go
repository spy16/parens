package parens

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ParseStr is a convenience wrapper for Parse.
func ParseStr(src string) (Expr, error) {
	return Parse(strings.NewReader(src))
}

// Parse parses till the EOF and returns all s-exprs as a single ModuleExpr.
// This should be used to build an entire module from a file or string etc.
func Parse(sc io.RuneScanner) (Expr, error) {
	me := ModuleExpr{}

	var expr Expr
	var err error
	for {
		expr, err = ParseOne(sc)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		me = append(me, expr)
	}

	return me, nil
}

// ParseOne consumes runes from the reader until a single s-expression is extracted.
// Returns any other errors from reader. This should be used when a continuous parse
// eval from a stream is necessary (e.g. TCP socket).
func ParseOne(sc io.RuneScanner) (Expr, error) {
	var expr Expr
	var err error
	for {
		expr, err = buildExpr(sc)
		if err != nil {
			return nil, err
		}

		// if expr is not nil, return it. otherwise, continue this loop
		// (for example, a whitespace can lead to nil)
		if expr != nil {
			return expr, nil
		}
	}
}

// Expr represents an expression.
type Expr interface {
	Eval(env Scope) (interface{}, error)
}

// Scope is responsible for managing bindings.
type Scope interface {
	Get(name string) (interface{}, error)
	Bind(name string, v interface{}, doc ...string) error
	Root() Scope
}

func buildExpr(rd io.RuneScanner) (Expr, error) {
	ru, _, err := rd.ReadRune()
	if err != nil {
		return nil, err
	}

	if unicode.IsControl(ru) {
		return nil, nil
	}

	switch ru {
	case '"':
		rd.UnreadRune()
		return buildStrExpr(rd)
	case '(':
		rd.UnreadRune()
		return buildListExpr(rd)
	case '[':
		rd.UnreadRune()
		return buildVectorExpr(rd)
	case '\'':
		rd.UnreadRune()
		return buildQuoteExpr(rd)
	case ':':
		rd.UnreadRune()
		return buildKeywordExpr(rd)
	case ';':
		rd.UnreadRune()
		_, err := buildCommentExpr(rd)
		if err != nil {
			return nil, err
		}
		return nil, nil

	case ' ', '\t', '\n':
		return nil, nil
	case ')', ']':
		return nil, io.EOF
	default:
		if utf8.ValidRune(ru) {
			rd.UnreadRune()
			return buildSymbolOrNumberExpr(rd)
		}
	}

	return nil, fmt.Errorf("invalid character '%c'", ru)

}

func buildListExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, '('); err != nil {
		return nil, err
	}

	lst := []Expr{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			return nil, err
		}

		if ru == ')' {
			break
		}

		rd.UnreadRune()

		expr, err := buildExpr(rd)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			lst = append(lst, expr)
		}
	}

	return ListExpr(lst), nil
}
func buildKeywordExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, ':'); err != nil {
		return nil, err
	}

	kw := []rune{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if isSepratingChar(ru) {
			rd.UnreadRune()
			break
		}

		if oneOf(ru, '\\') {
			return nil, fmt.Errorf("unexpected character '%c'", ru)
		}

		kw = append(kw, ru)
	}

	return KeywordExpr(kw), nil
}

func buildCommentExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, ';'); err != nil {
		return nil, err
	}

	comment := []rune{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if ru == '\n' {
			break
		}

		comment = append(comment, ru)
	}

	return CommentExpr(strings.TrimSpace(string(comment))), nil
}

func buildQuoteExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, '\''); err != nil {
		return nil, err
	}

	expr, err := buildExpr(rd)
	if err != nil {
		return nil, err
	}

	return QuoteExpr{Expr: expr}, nil
}

// TODO:
// - Support for hex (0x) and binary (0x) numbers
// - Support for scientific notation (1.3e10)
// - Clear differentiation between symbol and number
func buildSymbolOrNumberExpr(rd io.RuneScanner) (Expr, error) {
	seq := []rune{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if isSepratingChar(ru) {
			rd.UnreadRune()
			break
		}

		seq = append(seq, ru)
	}

	s := string(seq)
	if numberRegex.MatchString(s) {
		return NumberExpr{
			NumStr: s,
		}, nil
	}

	return SymbolExpr(s), nil
}

func buildVectorExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, '['); err != nil {
		return nil, err
	}

	vals := []Expr{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			return nil, err
		}

		if ru == ']' {
			break
		}
		rd.UnreadRune()

		expr, err := buildExpr(rd)
		if err != nil {
			return nil, err
		}

		if expr != nil {
			vals = append(vals, expr)
		}
	}

	return VectorExpr(vals), nil
}

func buildStrExpr(rd io.RuneScanner) (Expr, error) {
	if err := ensurePrefix(rd, '"'); err != nil {
		return nil, err
	}

	val := []rune{}

	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			return nil, err
		}

		lastI := len(val) - 1
		if len(val) >= 1 && val[lastI] == '\\' {
			var esc byte
			switch ru {
			case 'n':
				esc = '\n'
			case 't':
				esc = '\t'
			case 'r':
				esc = '\r'
			case '"':
				esc = '"'
			}

			if esc != 0 {
				val[lastI] = rune(esc)
				continue
			}

		}
		if oneOf(ru, 't', 'n', '"', 'r') {
		}
		if ru == '"' {
			break
		}

		val = append(val, ru)
	}

	return StringExpr(val), nil
}

var numberRegex = regexp.MustCompile("^(\\+|-)?\\d+(\\.\\d+)?$")

func ensurePrefix(rd io.RuneScanner, prefix rune) error {
	ru, _, err := rd.ReadRune()
	if err != nil {
		return err
	}

	if ru != prefix {
		return fmt.Errorf("expected '%c' at the beginning, found '%c'", prefix, ru)
	}

	return nil
}

func isSepratingChar(ru rune) bool {
	return oneOf(ru, ' ', '\t', '\n', '\r', '(', ')', '[', ']', '{', '}', '"', '\'')
}

func oneOf(ru rune, set ...rune) bool {
	for _, rs := range set {
		if ru == rs {
			return true
		}
	}
	return false
}
