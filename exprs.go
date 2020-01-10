package parens

import (
	"fmt"
	"strings"
)

// MacroFunc represents the signature of the Go macro functions. Functions
// bound in the scope as MacroFunc will receive un-evaluated list of s-exps
// and the current scope.
type MacroFunc func(scope Scope, exprs []Expr) (interface{}, error)

// Float64 represents double precision floating point numbers represented
// using float or scientific number formats.
type Float64 float64

// Eval returns the underlying double-precision float value.
func (f64 Float64) Eval(scope Scope) (interface{}, error) {
	return f64, nil
}

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Eval returns the underlying integer value.
func (i64 Int64) Eval(scope Scope) (interface{}, error) {
	return i64, nil
}

func (i64 Int64) String() string { return fmt.Sprintf("%d", i64) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Eval returns the unquoted string value.
func (se String) Eval(scope Scope) (interface{}, error) { return se, nil }

func (se String) String() string { return fmt.Sprintf("\"%s\"", string(se)) }

// Character represents a character literal.  For example, \a, \b, \1, \∂ etc
// are valid character literals. In addition, special literals like \newline,
// \space etc are supported.
type Character rune

// Eval returns the character value.
func (char Character) Eval(scope Scope) (interface{}, error) { return char, nil }

func (char Character) String() string { return fmt.Sprintf("\\%c", rune(char)) }

// Keyword represents a keyword literal.
type Keyword string

// Eval returns the keyword value.
func (kw Keyword) Eval(scope Scope) (interface{}, error) { return kw, nil }

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Symbol represents a name given to a value in memory.
type Symbol string

// Eval returns the value bound for the symbol in the scope.
func (sym Symbol) Eval(scope Scope) (interface{}, error) {
	return scope.Get(string(sym))
}

func (sym Symbol) String() string { return string(sym) }

// List represents an list of forms. Evaluating a list leads to a function
// invocation.
type List []Expr

// Eval executes the list as a function invocation.
func (lf List) Eval(scope Scope) (interface{}, error) {
	if len(lf) == 0 {
		return lf, nil
	}

	val, err := lf[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if macroFn, ok := val.(MacroFunc); ok {
		return macroFn(scope, lf[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(lf); i++ {
		arg, err := lf[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return reflectCall(val, args...)
}

func (lf List) String() string { return containerString(lf, "(", ")", " ") }

// Vector represents a list of values. Unlike List type, evaluation of
// vector does not lead to function invoke.
type Vector []Expr

// Eval evaluates each item in the vector and returns the result list.
func (vf Vector) Eval(scope Scope) (interface{}, error) { return evalForms(scope, vf) }

func (vf Vector) String() string { return containerString(vf, "[", "]", " ") }

// Module represents a group of forms. Evaluating a module form returns the
// result of evaluating the last form in the list.
type Module []Expr

// Eval evaluates all the forms and returns the result of the last evaluation.
func (m Module) Eval(scope Scope) (interface{}, error) {
	if len(m) == 0 {
		return nil, nil
	}

	val, err := evalForms(scope, m)
	if err != nil {
		return nil, err
	}

	return val[len(val)-1], nil
}

func (m Module) String() string { return containerString(m, "", "", "\n") }

func containerString(forms []Expr, begin, end, sep string) string {
	parts := make([]string, len(forms))
	for i, expr := range forms {
		parts[i] = fmt.Sprintf("%v", expr)
	}
	return begin + strings.Join(parts, sep) + end
}

func evalForms(scope Scope, forms []Expr) ([]interface{}, error) {
	var res []interface{}

	for _, form := range forms {
		val, err := form.Eval(scope)
		if err != nil {
			return nil, err
		}

		res = append(res, val)
	}

	return res, nil
}
