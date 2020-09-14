package value

import (
	"fmt"
	"math"
	"sync"
)

// Nil represents the Value 'nil'.
type Nil struct{}

// SExpr returns a valid s-expression representing Nil.
func (Nil) SExpr() (string, error) { return "nil", nil }

// Int64 represents a 64-bit integer Value.
type Int64 int64

// SExpr returns a valid s-expression representing Int64.
func (i64 Int64) SExpr() (string, error) { return fmt.Sprintf("%d", i64), nil }

// Equals returns true if the other Value is also an integer and has same Value.
func (i64 Int64) Equals(other Any) bool {
	val, isInt := other.(Int64)
	return isInt && (val == i64)
}

// Float64 represents a 64-bit double precision floating point Value.
type Float64 float64

// SExpr returns a valid s-expression representing Float64.
func (f64 Float64) SExpr() (string, error) {
	if math.Abs(float64(f64)) >= 1e16 {
		return fmt.Sprintf("%e", f64), nil
	}

	return fmt.Sprintf("%f", f64), nil
}

// Equals returns true if 'other' is also a float and has same Value.
func (f64 Float64) Equals(other Any) bool {
	val, isFloat := other.(Float64)
	return isFloat && (val == f64)
}

// Bool represents a boolean Value.
type Bool bool

// SExpr returns a valid s-expression representing Bool.
func (b Bool) SExpr() (string, error) {
	if bool(b) {
		return "true", nil
	}

	return "false", nil
}

// Equals returns true if 'other' is a boolean and has same logical Value.
func (b Bool) Equals(other Any) bool {
	val, ok := other.(Bool)
	return ok && (val == b)
}

// Char represents a Unicode character.
type Char rune

// SExpr returns a valid s-expression representing Char.
func (char Char) SExpr() (string, error) {
	return fmt.Sprintf("\\%c", char), nil
}

// Equals returns true if the other Value is also a character and has same Value.
func (char Char) Equals(other Any) bool {
	val, isChar := other.(Char)
	return isChar && (val == char)
}

// String represents a string of characters.
type String struct {
	Value string
	hash  int64
	once  sync.Once
}

// SExpr returns a valid s-expression representing String.
func (str *String) SExpr() (string, error) {
	return fmt.Sprintf("\"%s\"", str.Value), nil
}

// Equals returns true if 'other' is string and has same Value.
func (str *String) Equals(other Any) bool {
	otherStr, isStr := other.(*String)
	return isStr && (otherStr.Value == str.Value)
}

// Symbol represents a lisp symbol Value.
type Symbol struct {
	Value string
	hash  int64
	once  sync.Once
}

// SExpr returns a valid s-expression representing Symbol.
func (sym *Symbol) SExpr() (string, error) { return sym.Value, nil }

// Equals returns true if the other Value is also a symbol and has same Value.
func (sym *Symbol) Equals(other Any) bool {
	otherSym, isSym := other.(*Symbol)
	return isSym && (sym.Value == otherSym.Value)
}

// Keyword represents a keyword Value.
type Keyword struct {
	Value string
	hash  int64
	once  sync.Once
}

// SExpr returns a valid s-expression representing Keyword.
func (kw *Keyword) SExpr() (string, error) { return fmt.Sprintf(":%s", kw.Value), nil }

// Equals returns true if the other Value is keyword and has same Value.
func (kw *Keyword) Equals(other Any) bool {
	otherKW, isKeyword := other.(*Keyword)
	return isKeyword && (otherKW == kw)
}
