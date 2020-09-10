package parens

// Option can be used with New() to customize initialization of VM.
type Option func(ev *Evaluator)

// WithMaxDepth sets the maximum stack depth allowed for invocations.
func WithMaxDepth(depth int) Option {
	return func(ev *Evaluator) { ev.maxDepth = depth }
}

// WithExpander sets the macro Expander to be used by the VM. If nil, a builtin
// Expander will be used.
func WithExpander(expander Expander) Option {
	return func(ev *Evaluator) {
		if expander == nil {
			expander = &builtinExpander{}
		}
		ev.expander = expander
	}
}

// WithAnalyzer sets the Analyzer to be used by the VM. If replaceBuiltin is set,
// given analyzer will be used for all forms. Otherwise, it will be used only for
// forms not handled by the builtinAnalyzer.
func WithAnalyzer(replaceBuiltin bool, analyzer Analyzer) Option {
	return func(ev *Evaluator) {
		switch {
		case analyzer == nil:
			ev.analyzer = &builtinAnalyzer{ev: ev}

		case replaceBuiltin:
			ev.analyzer = analyzer

		default:
			ev.analyzer = &builtinAnalyzer{ev: ev, extender: analyzer}
		}
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAnalyzer(true, nil),
		WithExpander(nil),
	}, opts...)
}
