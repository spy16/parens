package stdlib

import (
	"fmt"
	"math"

	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// WithBuiltins registers different built-in functions into the
// given scope.
func WithBuiltins(scope *reflection.Scope) {
	builtins := map[string]interface{}{
		// macros
		"begin":   parser.MacroFunc(Begin),
		"setq":    parser.MacroFunc(Setq),
		"cond":    parser.MacroFunc(Conditional),
		"let":     parser.MacroFunc(Let),
		"inspect": parser.MacroFunc(Inspect),

		// functions
		"print":   fmt.Print,
		"println": fmt.Println,
		"printf":  fmt.Printf,
		"+":       Add,
		"-":       Sub,
		"*":       Mul,
		"/":       Div,
		"^":       math.Pow,

		// values
		"true":  true,
		"false": false,
		"nil":   false,
	}

	for name, val := range builtins {
		scope.Bind(name, val)
	}
}
