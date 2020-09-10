package parens

import (
	"github.com/spy16/parens/value"
)

// New returns a new instance of Evaluator with empty Global frame.
func New(opts ...Option) *Evaluator {
	var ev Evaluator
	for _, opt := range withDefaults(opts) {
		opt(&ev)
	}
	return &ev
}

// Evaluator provides the environment where forms are evaluated for result. Forms are
// first converted to Expr based on pre-defined rules or custom analyzer rules
// (See Analyzer). Evaluator is safe for concurrent use.
type Evaluator struct {
	analyzer Analyzer
	expander Expander
}

// Expr represents an expression that can be evaluated against the Evaluator.
type Expr interface {
	Eval(ctx Context, ev *Evaluator) (value.Any, error)
}

// Invokable represents a value that can be invoked for result.
type Invokable interface {
	Invoke(ev *Evaluator, args ...value.Any) (value.Any, error)
}

// Analyzer implementation can be set on the Evaluator to override evaluation
// rules for forms. See WithAnalyzer().
type Analyzer interface {
	Analyze(form value.Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(ev *Evaluator, form value.Any) (value.Any, error)
}

// Eval performs macro-expansion if necessary, converts the expanded form to an
// expression and evaluates the resulting expression.
func (ev *Evaluator) Eval(ctx Context, form value.Any) (value.Any, error) {
	if form == nil {
		return value.Nil{}, nil
	}

	expr, err := ev.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return value.Nil{}, nil
	}

	return expr.Eval(ctx, ev)
}

func (ev *Evaluator) expandAnalyze(form value.Any) (Expr, error) {
	if expr, ok := form.(Expr); ok {
		// Already an Expr, nothing to do.
		return expr, nil
	}

	if expanded, err := ev.expander.Expand(ev, form); err != nil {
		return nil, err
	} else if expanded != nil {
		// Expansion did happen. Throw away the old form and continue with
		// the expanded version.
		form = expanded
	}

	return ev.analyzer.Analyze(form)
}
