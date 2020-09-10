package parens

// Option can be used with New() to customize initialization of Evaluator
// Instance.
type Option func(p *Evaluator)

// WithExpander sets the macro Expander to be used by the p. If nil, a builtin
// Expander will be used.
func WithExpander(expander Expander) Option {
	return func(p *Evaluator) {
		if expander == nil {
			expander = &builtinExpander{}
		}
		p.expander = expander
	}
}

// WithAnalyzer sets the Analyzer to be used by the p. If replaceBuiltin is set,
// given analyzer will be used for all forms. Otherwise, it will be used only for
// forms not handled by the builtinAnalyzer.
func WithAnalyzer(replaceBuiltin bool, analyzer Analyzer) Option {
	return func(p *Evaluator) {
		switch {
		case analyzer == nil:
			p.analyzer = &builtinAnalyzer{}

		case replaceBuiltin:
			p.analyzer = analyzer

		default:
			p.analyzer = &builtinAnalyzer{extender: analyzer}
		}
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAnalyzer(true, nil),
		WithExpander(nil),
	}, opts...)
}
