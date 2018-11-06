package stdlib

import (
	"fmt"
	"math"

	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// RegisterBuiltins registers different built-in functions into the
// given scope.
func RegisterBuiltins(scope *reflection.Scope) {
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
		scope.Bind(name, val)
	}
}
