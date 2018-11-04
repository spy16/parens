package stdlib

import (
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// WithBuiltins registers different built-in functions into the
// given scope.
func WithBuiltins(scope *reflection.Scope) {
	builtins := map[string]interface{}{
		"setq":    parser.MacroFunc(Setq),
		"inspect": parser.MacroFunc(Inspect),
	}

	for name, val := range builtins {
		scope.Bind(name, val)
	}
}
