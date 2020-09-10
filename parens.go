package parens

import (
	"errors"

	"github.com/spy16/parens/value"
)

const globalFrame = "<global>"

// ErrNotFound is returned by Evaluator when a symbol resolution fails.
var ErrNotFound = errors.New("not found")

// New returns a new instance of Evaluator with empty Global frame.
func New(opts ...Option) *Evaluator {
	var ev Evaluator
	ev.push(stackFrame{name: globalFrame})
	for _, opt := range withDefaults(opts) {
		opt(&ev)
	}
	return &ev
}

// Evaluator provides the environment where forms are evaluated for result. Forms are
// first converted to Expr based on pre-defined rules or custom analyzer rules
// (See WithAnalyzer). Evaluator is not safe for concurrent use without external sync.
type Evaluator struct {
	stack    []stackFrame
	maxDepth int
	analyzer Analyzer
	expander Expander
}

// Expr represents an expression that can be evaluated against the Evaluator.
type Expr interface {
	Eval(*Evaluator) (value.Any, error)
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
func (ev *Evaluator) Eval(form value.Any) (value.Any, error) {
	if form == nil {
		return value.Nil{}, nil
	}

	expr, err := ev.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return value.Nil{}, nil
	}

	return expr.Eval(ev)
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

func (ev *Evaluator) resolve(name string) (value.Any, error) {
	if len(ev.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	// traverse from top of the stack until a binding is found.
	for i := len(ev.stack) - 1; i >= 0; i-- {
		if v, found := ev.stack[i].get(name); found {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

func (ev *Evaluator) push(frame stackFrame) {
	ev.stack = append(ev.stack, frame)
}

func (ev *Evaluator) pop() *stackFrame {
	if len(ev.stack) == 0 {
		panic("Evaluator stack must never be empty")
	}

	f := ev.stack[len(ev.stack)-1]
	ev.stack = ev.stack[0 : len(ev.stack)-1]
	return &f
}

type stackFrame struct {
	name string
	args []value.Any
	vars map[string]value.Any

	// positional information
	file      string
	line, col int
}

func (frame *stackFrame) get(name string) (value.Any, bool) {
	val, found := frame.vars[name]
	return val, found
}

func (frame *stackFrame) set(name string, val value.Any) {
	if frame.vars == nil {
		frame.vars = map[string]value.Any{}
	}
	frame.vars[name] = val
}
