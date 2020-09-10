package parens

// New returns a new instance of parens initialised based on given
// options.
func New(opts ...Option) *Evaluator {
	p := &Evaluator{}
	for _, opt := range withDefaults(opts) {
		opt(p)
	}
	return p
}

// Any represents any value.
type Any interface{}

// Evaluator represents an instance of parens interpreter.
type Evaluator struct {
	expander Expander
	analyzer Analyzer
}

// Eval performs macro-expansion if necessary, converts the expanded form
// to an expression and evaluates the resulting expression.
func (e *Evaluator) Eval(form Any) (Any, error) {
	expr, err := e.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return nil, nil
	}

	return expr.Eval()
}

func (e *Evaluator) expandAnalyze(form Any) (Expr, error) {
	if expr, ok := form.(Expr); ok {
		// Already an Expr, nothing to do.
		return expr, nil
	}

	if expanded, err := e.expander.Expand(e, form); err != nil {
		return nil, err
	} else if expanded != nil {
		// Expansion did happen. Throw away the old form and continue with
		// the expanded version.
		form = expanded
	}

	return e.analyzer.Analyze(form)
}

// Analyzer implementation is responsible for performing syntax analysis
// on given form.
type Analyzer interface {
	// Analyze should perform syntax checks for special forms etc. and
	// return Expr values that can be evaluated against a context.
	Analyze(form Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(p *Evaluator, form Any) (Any, error)
}

// Expr represents an expression that can be evaluated against a context.
type Expr interface {
	// TODO: (spy16) Modify signature as per Expr implementations.
	Eval() (Any, error)
}
