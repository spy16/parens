package reader

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spy16/parens"
)

// Macro implementations can be plugged into the Reader to extend, override
// or customize behavior of the reader.
type Macro func(rd *Reader, init rune) (parens.Any, error)

// // TODO(enhancement):  implement parens.Set
// // SetReader implements the reader macro for reading set from source.
// func SetReader(setEnd rune, factory func() parens.Set) Macro {
// 	return func(rd *Reader, _ rune) (parens.Any, error) {
// 		forms, err := rd.Container(setEnd, "Set")
// 		if err != nil {
// 			return nil, err
// 		}
// 		return factory().Conj(forms...), nil
// 	}
// }

// // TODO(enhancement): implement parens.Vector
// // VectorReader implements the reader macro for reading vector from source.
// func VectorReader(vecEnd rune, factory func() parens.Vector) Macro {
// 	return func(rd *Reader, _ rune) (parens.Any, error) {
// 		forms, err := rd.Container(vecEnd, "Vector")
// 		if err != nil {
// 			return nil, err
// 		}

// 		vec := factory()
// 		for _, f := range forms {
// 			_ = f
// 		}

// 		return vec, nil
// 	}
// }

// // TODO(enhancement) implement parens.Map
// // MapReader returns a reader macro for reading map values from source. factory
// // is used to construct the map and `Assoc` is called for every pair read.
// func MapReader(mapEnd rune, factory func() parens.Map) Macro {
// 	return func(rd *Reader, _ rune) (parens.Any, error) {
// 		forms, err := rd.Container(mapEnd, "Map")
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(forms)%2 != 0 {
// 			return nil, errors.New("expecting even number of forms within {}")
// 		}

// 		m := factory()
// 		for i := 0; i < len(forms); i += 2 {
// 			if m.HasKey(forms[i]) {
// 				return nil, fmt.Errorf("duplicate key: %v", forms[i])
// 			}

// 			m, err = m.Assoc(forms[i], forms[i+1])
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		return m, nil
// 	}
// }

// UnmatchedDelimiter implements a reader macro that can be used to capture
// unmatched delimiters such as closing parenthesis etc.
func UnmatchedDelimiter() Macro {
	return func(_ *Reader, initRune rune) (parens.Any, error) {
		return nil, Error{
			Cause: ErrUnmatchedDelimiter,
			Rune:  initRune,
		}
	}
}

func readNumber(rd *Reader, init rune) (parens.Any, error) {
	numStr, err := rd.Token(init)
	if err != nil {
		return nil, err
	}

	decimalPoint := strings.ContainsRune(numStr, '.')
	isRadix := strings.ContainsRune(numStr, 'r')
	isScientific := strings.ContainsRune(numStr, 'e')

	switch {
	case isRadix && (decimalPoint || isScientific):
		return nil, Error{
			Cause: ErrNumberFormat,
			Form:  numStr,
		}

	case isScientific:
		return parseScientific(numStr)

	case decimalPoint:
		v, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, Error{
				Cause: ErrNumberFormat,
				Form:  numStr,
			}
		}
		return parens.Float64(v), nil

	case isRadix:
		return parseRadix(numStr)

	default:
		v, err := strconv.ParseInt(numStr, 0, 64)
		if err != nil {
			return nil, Error{
				Cause: ErrNumberFormat,
				Form:  numStr,
			}
		}

		return parens.Int64(v), nil
	}
}

func readSymbol(rd *Reader, init rune) (parens.Any, error) {
	s, err := rd.Token(init)
	if err != nil {
		return nil, err
	}

	return parens.Symbol(s), nil
}

func readString(rd *Reader, _ rune) (parens.Any, error) {
	var b strings.Builder

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, Error{
					Form:  "string",
					Cause: ErrEOF,
				}
			}

			return nil, err
		}

		if r == '\\' {
			r2, err := rd.NextRune()
			if err != nil {
				if err == io.EOF {
					return nil, Error{
						Form:  "string",
						Cause: ErrEOF,
					}
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

	return parens.String(b.String()), nil
}

func readComment(rd *Reader, _ rune) (parens.Any, error) {
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

func readKeyword(rd *Reader, init rune) (parens.Any, error) {
	token, err := rd.Token(-1)
	if err != nil {
		return nil, err
	}

	return parens.Keyword(token), nil
}

func readCharacter(rd *Reader, _ rune) (parens.Any, error) {
	r, err := rd.NextRune()
	if err != nil {
		return nil, Error{
			Form:  "character",
			Cause: ErrEOF,
		}
	}

	token, err := rd.Token(r)
	if err != nil {
		return nil, err
	}
	runes := []rune(token)

	if len(runes) == 1 {
		return parens.Char(runes[0]), nil
	}

	v, found := charLiterals[token]
	if found {
		return parens.Char(v), nil
	}

	if token[0] == 'u' {
		return readUnicodeChar(token[1:], 16)
	}

	return nil, fmt.Errorf("unsupported character: '\\%s'", token)
}

func readList(rd *Reader, _ rune) (parens.Any, error) {
	const listEnd = ')'

	forms := make([]parens.Any, 0, 32) // pre-allocate to improve performance on small lists
	if err := rd.Container(listEnd, "list", func(val parens.Any) error {
		forms = append(forms, val)
		return nil
	}); err != nil {
		return nil, err
	}

	return parens.NewList(forms...), nil
}

func quoteFormReader(expandFunc string) Macro {
	return func(rd *Reader, _ rune) (parens.Any, error) {
		expr, err := rd.One()
		if err != nil {
			if err == io.EOF {
				return nil, Error{
					Form:  expandFunc,
					Cause: ErrEOF,
				}
			} else if err == ErrSkip {
				return nil, Error{
					Form:  expandFunc,
					Cause: errors.New("cannot quote a no-op form"),
				}
			}
			return nil, err
		}

		return parens.NewList(parens.Symbol(expandFunc), expr), nil
	}
}
