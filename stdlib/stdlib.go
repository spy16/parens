package stdlib

import (
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// WithBuiltins registers different built-in functions into the
// given scope.
func WithBuiltins(scope *reflection.Scope) {
	builtins := map[string]interface{}{
		// macros
		"setq": parser.MacroFunc(Setq),
		"cond": parser.MacroFunc(Conditional),
		"let":  parser.MacroFunc(Let),

		// functions
		"inspect": parser.MacroFunc(Inspect),

		// values
		"true":  true,
		"false": false,
		"nil":   false,
	}

	for name, val := range builtins {
		scope.Bind(name, val)
	}
}
