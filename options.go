package parens

// Option can be used with New() to customize initialization of Evaluator
// Instance.
type Option func(ctx *Context)

// WithMapFactory sets the factory to be used for creating variables map
// in a stack frame.
func WithMapFactory(factory func() ConcurrentMap) Option {
	return func(ctx *Context) {
		if factory == nil {
			factory = func() ConcurrentMap { return &mutexMap{} }
		}
		ctx.mapFactory = factory
	}
}

// WithMaxDepth sets the max depth allowed for stack.
func WithMaxDepth(depth int) Option {
	return func(ctx *Context) {
		ctx.maxDepth = depth
	}
}

// WithExpander sets the macro Expander to be used by the p. If nil, a builtin
// Expander will be used.
func WithExpander(expander Expander) Option {
	return func(env *Context) {
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
	return func(ctx *Context) {
		if analyzer == nil {
			analyzer = &BasicAnalyzer{
				SpecialForms: map[string]ParseSpecial{
					"go":    parseGoExpr,
					"def":   parseDefExpr,
					"quote": parseQuoteExpr,
				},
			}
		}
		ctx.analyzer = analyzer
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAnalyzer(nil),
		WithExpander(nil),
		WithMapFactory(nil),
	}, opts...)
}
