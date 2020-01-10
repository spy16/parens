package parens

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

var (
	// ErrSkip can be returned by reader macro to indicate a no-op.
	ErrSkip      = errors.New("skip expr")
	errStringEOF = errors.New("EOF while reading string")
	errCharEOF   = errors.New("EOF while reading character")
)

var (
	escapeMap = map[rune]rune{
		'"':  '"',
		'n':  '\n',
		'\\': '\\',
		't':  '\t',
		'a':  '\a',
		'f':  '\a',
		'r':  '\r',
		'b':  '\b',
		'v':  '\v',
	}

	charLiterals = map[string]rune{
		"tab":       '\t',
		"space":     ' ',
		"newline":   '\n',
		"return":    '\r',
		"backspace": '\b',
		"formfeed":  '\f',
	}
)

// New returns a lisp reader instance which can read forms from r. Reader
// behavior can be customized by using SetMacro to override or remove from
// the default read table. File name will be inferred from the reader value
// and type information.
func New(rs io.Reader) *Reader {
	rd := &Reader{
		Stream: Stream{
			File: inferFileName(rs),
			rs:   bufio.NewReader(rs),
		},
		macros: defaultReadTable(),
	}

	return rd
}

// ReaderMacro implementations can be plugged into the Reader to extend, override
// or customize behavior of the reader.
type ReaderMacro func(rd *Reader, init rune) (Expr, error)

// ParseStr is a convenience wrapper for Parse.
func ParseStr(src string) (Expr, error) {
	return Parse(strings.NewReader(src))
}

// Parse parses till the EOF and returns all s-exprs as a single ModuleExpr.
// This should be used to build an entire module from a file or string etc.
func Parse(r io.Reader) (Expr, error) {
	return New(r).All()
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

// Reader provides functions to parse characters from a stream into symbolic
// expressions or forms.
type Reader struct {
	Stream

	Hook   ReaderMacro
	macros map[rune]ReaderMacro
}

// All consumes characters from stream until EOF and returns a list of all the
// forms parsed. Any no-op forms (e.g., comment) returned will not be included
// in the result.
func (rd *Reader) All() (Module, error) {
	var forms []Expr

	for {
		form, err := rd.One()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

// One consumes characters from underlying stream until a complete form is
// parsed and returns the form while ignoring the no-op forms like comments.
// Except EOF, all errors will be wrapped with ReaderError type along with
// the positional information obtained using Info().
func (rd *Reader) One() (Expr, error) {
	for {
		form, err := rd.readOne()
		if err != nil {
			if err == ErrSkip {
				continue
			}

			return nil, rd.annotateErr(err)
		}

		return form, nil
	}
}

// IsTerminal returns true if the rune should terminate a form. ReaderMacro
// trigger runes defined in the read table and all space characters including
// "," are considered terminal.
func (rd *Reader) IsTerminal(r rune) bool {
	_, found := rd.macros[r]
	return found || isSpace(r)
}

// SetMacro sets the given reader macro as the handler for init rune in the
// read table. Overwrites if a macro is already present. If the macro value
// given is nil, entry for the init rune will be removed from the read table.
func (rd *Reader) SetMacro(init rune, macro ReaderMacro) {
	if macro == nil {
		delete(rd.macros, init)
		return
	}

	rd.macros[init] = macro
}

// readOne is same as One() but always returns un-annotated errors.
func (rd *Reader) readOne() (Expr, error) {
	if err := rd.SkipSpaces(); err != nil {
		return nil, err
	}

	r, err := rd.NextRune()
	if err != nil {
		return nil, err
	}

	if unicode.IsNumber(r) {
		return readNumber(rd, r)
	} else if r == '+' || r == '-' {
		r2, err := rd.NextRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != io.EOF {
			rd.Unread(r2)
			if unicode.IsNumber(r2) {
				return readNumber(rd, r)
			}
		}
	}

	macro, found := rd.macros[r]
	if found {
		return macro(rd, r)
	}

	if rd.Hook != nil {
		f, err := rd.Hook(rd, r)
		if err != ErrSkip {
			return f, err
		}
	}

	return readSymbol(rd, r)
}

func (rd *Reader) annotateErr(e error) error {
	if e == io.EOF || e == ErrSkip {
		return e
	}

	file, line, col := rd.Info()
	return ReaderError{
		Cause:  e,
		File:   file,
		Line:   line,
		Column: col,
	}
}

func readNumber(rd *Reader, init rune) (Expr, error) {
	numStr, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	decimalPoint := strings.ContainsRune(numStr, '.')
	isRadix := strings.ContainsRune(numStr, 'r')
	isScientific := strings.ContainsRune(numStr, 'e')

	switch {
	case isRadix && (decimalPoint || isScientific):
		return nil, fmt.Errorf("illegal number format: '%s'", numStr)

	case isScientific:
		return parseScientific(numStr)

	case decimalPoint:
		v, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, errors.Wrap(err, "illegal number format")
			// return nil, fmt.Errorf("illegal number format: '%s'", numStr)
		}
		return Float64(v), nil

	case isRadix:
		return parseRadix(numStr)

	default:
		v, err := strconv.ParseInt(numStr, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "illegal number format")
			// return nil, fmt.Errorf("illegal number format '%s'", numStr)
		}

		return Int64(v), nil
	}
}

func readSymbol(rd *Reader, init rune) (Expr, error) {
	s, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return Symbol(s), nil
}

func readKeyword(rd *Reader, _ rune) (Expr, error) {
	r, err := rd.NextRune()
	if err != nil {
		return nil, err
	}

	token, err := readToken(rd, r)
	if err != nil {
		return nil, err
	}

	return Keyword(token), nil
}

func readCharacter(rd *Reader, _ rune) (Expr, error) {
	r, err := rd.NextRune()
	if err != nil {
		return nil, errCharEOF
	}

	token, err := readToken(rd, r)
	if err != nil {
		return nil, err
	}
	runes := []rune(token)

	if len(runes) == 1 {
		return Character(runes[0]), nil
	}

	v, found := charLiterals[token]
	if found {
		return Character(v), nil
	}

	if token[0] == 'u' {
		return readUnicodeChar(token[1:], 16)
	}

	return nil, fmt.Errorf("unsupported character: '\\%s'", token)
}

func readUnicodeChar(token string, base int) (Character, error) {
	num, err := strconv.ParseInt(token, base, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	if num < 0 || num >= unicode.MaxRune {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	return Character(num), nil
}

func readString(rd *Reader, _ rune) (Expr, error) {
	var b strings.Builder

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, errStringEOF
			}

			return nil, err
		}

		if r == '\\' {
			r2, err := rd.NextRune()
			if err != nil {
				if err == io.EOF {
					return nil, errStringEOF
				}

				return nil, err
			}

			// TODO: Support for Unicode escape \uNN format.

			escaped, err := getEscape(r2)
			if err != nil {
				return nil, err
			}
			r = escaped
		} else if r == '"' {
			break
		}

		b.WriteRune(r)
	}

	return String(b.String()), nil
}

func readList(rd *Reader, _ rune) (Expr, error) {
	forms, err := readContainer(rd, '(', ')', "list")
	if err != nil {
		return nil, err
	}

	return List(forms), nil
}

func readVector(rd *Reader, _ rune) (Expr, error) {
	forms, err := readContainer(rd, '[', ']', "vector")
	if err != nil {
		return nil, err
	}

	return Vector(forms), nil
}

func readComment(rd *Reader, _ rune) (Expr, error) {
	for {
		r, err := rd.NextRune()
		if err != nil {
			return nil, err
		}

		if r == '\n' {
			break
		}
	}

	return nil, ErrSkip
}

func quoteFormReader(expandFunc string) ReaderMacro {
	return func(rd *Reader, _ rune) (Expr, error) {
		expr, err := rd.One()
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("EOF while reading quote form")
			} else if err == ErrSkip {
				return nil, errors.New("no-op form while reading quote form")
			}
			return nil, err
		}

		return List{
			Symbol(expandFunc),
			expr,
		}, nil
	}
}

func unmatchedDelimiter(_ *Reader, initRune rune) (Expr, error) {
	return nil, fmt.Errorf("unmatched delimiter '%c'", initRune)
}

func readToken(rd *Reader, init rune) (string, error) {
	var b strings.Builder
	if init != -1 {
		b.WriteRune(init)
	}

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if rd.IsTerminal(r) {
			rd.Unread(r)
			break
		}

		b.WriteRune(r)
	}

	return b.String(), nil
}

func readContainer(rd *Reader, _ rune, end rune, formType string) ([]Expr, error) {
	var forms []Expr

	for {
		if err := rd.SkipSpaces(); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("EOF while reading %s", formType)
			}
			return nil, err
		}

		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("EOF while reading %s", formType)
			}
			return nil, err
		}

		if r == end {
			break
		}
		rd.Unread(r)

		expr, err := rd.readOne()
		if err != nil {
			if err == ErrSkip {
				continue
			}
			return nil, err
		}
		forms = append(forms, expr)
	}

	return forms, nil
}

func parseRadix(numStr string) (Int64, error) {
	parts := strings.Split(numStr, "r")
	if len(parts) != 2 {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	base, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	repr := parts[1]
	if base < 0 {
		base = -1 * base
		repr = "-" + repr
	}

	v, err := strconv.ParseInt(repr, int(base), 64)
	if err != nil {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	return Int64(v), nil
}

func parseScientific(numStr string) (Float64, error) {
	parts := strings.Split(numStr, "e")
	if len(parts) != 2 {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	base, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	pow, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	return Float64(base * math.Pow(10, float64(pow))), nil
}

func getEscape(r rune) (rune, error) {
	escaped, found := escapeMap[r]
	if !found {
		return -1, fmt.Errorf("illegal escape sequence '\\%c'", r)
	}

	return escaped, nil
}

func inferFileName(rs io.Reader) string {
	switch r := rs.(type) {
	case *os.File:
		return r.Name()

	case *strings.Reader:
		return "<string>"

	case *bytes.Reader:
		return "<bytes>"

	case net.Conn:
		return fmt.Sprintf("<con:%s>", r.LocalAddr())

	default:
		return fmt.Sprintf("<%s>", reflect.TypeOf(rs))
	}
}

func defaultReadTable() map[rune]ReaderMacro {
	return map[rune]ReaderMacro{
		'"':  readString,
		';':  readComment,
		':':  readKeyword,
		'\\': readCharacter,
		'\'': quoteFormReader("quote"),
		'~':  quoteFormReader("unquote"),
		'(':  readList,
		')':  unmatchedDelimiter,
		'[':  readVector,
		']':  unmatchedDelimiter,
	}
}

// ReaderError wraps the parsing error with file and positional information.
type ReaderError struct {
	Cause  error
	File   string
	Line   int
	Column int
}

// Unwrap returns the error's cause
func (err ReaderError) Unwrap() error {
	return err.Cause
}

func (err ReaderError) Error() string {
	if e, ok := err.Cause.(ReaderError); ok {
		return e.Error()
	}

	return fmt.Sprintf("syntax error in '%s' (Line %d Col %d): %v", err.File, err.Line, err.Column, err.Cause)
}
