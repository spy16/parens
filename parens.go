package parens

import "github.com/spy16/parens/value"

// Evaluator provides the environment where forms are evaluated for result. Forms are
// first converted to Expr based on pre-defined rules or custom analyzer rules
// (See Analyzer).
//
// The zero value is ready to use.  Evaluator is safe for concurrent use.
type Evaluator struct {
	Analyzer
	Expander
}

// Expr represents an expression that can be evaluated against the Evaluator.
type Expr interface {
	Eval(ctx Context, ev Evaluator) (value.Any, error)
}

// Invokable represents a value that can be invoked for result.
type Invokable interface {
	Invoke(ev Evaluator, args ...value.Any) (value.Any, error)
}

// Analyzer implementation can be set on the Evaluator to override evaluation
// rules for forms. See WithAnalyzer().
type Analyzer interface {
	Analyze(ev Evaluator, form value.Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(ev Evaluator, form value.Any) (value.Any, error)
}

// Eval performs macro-expansion if necessary, converts the expanded form to an
// expression and evaluates the resulting expression.
func (ev Evaluator) Eval(ctx Context, form value.Any) (value.Any, error) {
	if form == nil {
		return nil, nil
	}

	expr, err := ev.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return nil, nil
	}

	return expr.Eval(ctx, ev)
}

// EvalAll evaluates each value in the list against the given env and returns a list
// of resultant value.
func (ev Evaluator) EvalAll(ctx Context, vals []value.Any) (res []value.Any, err error) {
	res = make([]value.Any, 0, len(vals))

	for _, form := range vals {
		if form, err = ev.Eval(ctx, form); err != nil {
			break
		}

		res = append(res, form)
	}

	return
}

func (ev Evaluator) expandAnalyze(form value.Any) (Expr, error) {
	if expr, ok := form.(Expr); ok {
		// Already an Expr, nothing to do.
		return expr, nil
	}

	if expanded, err := ev.expand(form); err != nil {
		return nil, err
	} else if expanded != nil {
		// Expansion did happen. Throw away the old form and continue with
		// the expanded version.
		form = expanded
	}

	return ev.analyze(form)
}

func (ev Evaluator) expand(form value.Any) (value.Any, error) {
	if ev.Expander == nil {
		return basicExpander{}.Expand(ev, form)
	}

	return ev.Expander.Expand(ev, form)
}

func (ev Evaluator) analyze(form value.Any) (Expr, error) {
	if ev.Analyzer == nil {
		return BasicAnalyzer{}.Analyze(ev, form)
	}

	return ev.Analyzer.Analyze(ev, form)
}
