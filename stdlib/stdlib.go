package stdlib

import (
	"fmt"
	"math"

	"github.com/spy16/parens/parser"
)

// RegisterBuiltins registers different built-in functions into the
// given scope.
func RegisterBuiltins(scope parser.Scope) error {
	builtins := map[string]interface{}{
		// macros
		"begin":   parser.MacroFunc(Begin),
		"setq":    parser.MacroFunc(Setq),
		"cond":    parser.MacroFunc(Conditional),
		"let":     parser.MacroFunc(Let),
		"inspect": parser.MacroFunc(Inspect),
		"lambda":  parser.MacroFunc(Lambda),

		// functions
		"print":   fmt.Print,
		"println": fmt.Println,
		"printf":  fmt.Printf,
		"+":       Add,
		"-":       Sub,
		"*":       Mul,
		"/":       Div,
		"^":       math.Pow,
		"==":      Eq,
		"not":     Not,
		">":       Gt,
		"<":       Lt,

		// values
		"true":  true,
		"false": false,
		"nil":   false,
	}

	for name, val := range builtins {
		if err := scope.Bind(name, val); err != nil {
			return err
		}
	}

	return nil
}
