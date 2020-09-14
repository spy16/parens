package parens

import "github.com/spy16/parens/value"

// Option can be used with New() to customize initialization of Evaluator
// Instance.
type Option func(env *Env)

// WithGlobals sets the global variables during initialisation.
func WithGlobals(globals map[string]value.Any) Option {
	return func(env *Env) {
		vars := env.mapFactory()
		for k, v := range globals {
			vars.Store(k, v)
		}
		env.stack[0].ConcurrentMap = vars
	}
}

// WithMapFactory sets the factory to be used for creating variables map
// in a stack frame.
func WithMapFactory(factory func() ConcurrentMap) Option {
	return func(env *Env) {
		if factory == nil {
			factory = newMutexMap
		} else {
			newMap := factory()
			for k, v := range env.stack[0].Map() {
				newMap.Store(k, v)
			}
			env.stack[0].ConcurrentMap = newMap
		}
		env.mapFactory = factory
	}
}

// WithMaxDepth sets the max depth allowed for stack.
func WithMaxDepth(depth int) Option {
	return func(env *Env) {
		env.maxDepth = depth
	}
}

// WithExpander sets the macro Expander to be used by the p. If nil, a builtin
// Expander will be used.
func WithExpander(expander Expander) Option {
	return func(env *Env) {
		if expander == nil {
			expander = &basicExpander{}
		}
		env.expander = expander
	}
}

// WithAnalyzer sets the Analyzer to be used by the p. If replaceBuiltin is set,
// given analyzer will be used for all forms. Otherwise, it will be used only for
// forms not handled by the builtinAnalyzer.
func WithAnalyzer(analyzer Analyzer) Option {
	return func(env *Env) {
		if analyzer == nil {
			analyzer = &BasicAnalyzer{
				SpecialForms: map[string]ParseSpecial{
					"go":    parseGoExpr,
					"def":   parseDefExpr,
					"quote": parseQuoteExpr,
				},
			}
		}
		env.analyzer = analyzer
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAnalyzer(nil),
		WithExpander(nil),
		WithMapFactory(nil),
	}, opts...)
}
