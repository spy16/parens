package parens

import (
	"context"

	"github.com/spy16/parens/value"
)

const globalFrame = "<global>"

// New returns a new root context initialised based on given options.
func New(opts ...Option) Env {
	env := Env{ctx: context.Background()}
	env.push(stackFrame{Name: globalFrame})
	for _, opt := range withDefaults(opts) {
		opt(&env)
	}
	return env
}

// Env represents the environment/context in which forms are evaluated
// for result. Env is not safe for concurrent use. Use fork() to get a
// child context for concurrent executions.
type Env struct {
	ctx        context.Context
	analyzer   Analyzer
	expander   Expander
	stack      []stackFrame
	maxDepth   int
	mapFactory func() ConcurrentMap
}

// Eval performs macro-expansion if necessary, converts the expanded form
// to an expression and evaluates the resulting expression.
func (env *Env) Eval(form value.Any) (value.Any, error) {
	expr, err := env.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return nil, nil
	}

	return expr.Eval(env)
}

func (env *Env) expandAnalyze(form value.Any) (Expr, error) {
	if expr, ok := form.(Expr); ok {
		// Already an Expr, nothing to do.
		return expr, nil
	}

	if expanded, err := env.expander.Expand(env, form); err != nil {
		return nil, err
	} else if expanded != nil {
		// Expansion did happen. Throw away the old form and continue with
		// the expanded version.
		form = expanded
	}

	return env.analyzer.Analyze(env, form)
}

// fork creates a child context from Env and returns it. The child context
// can be used as context for an independent thread of execution.
func (env *Env) fork() *Env {
	return &Env{
		ctx:        env.ctx,
		stack:      append([]stackFrame(nil), env.stack...),
		expander:   env.expander,
		analyzer:   env.analyzer,
		maxDepth:   env.maxDepth,
		mapFactory: env.mapFactory,
	}
}

func (env *Env) push(frame stackFrame) {
	env.stack = append(env.stack, frame)
}

func (env *Env) pop() (frame *stackFrame) {
	if len(env.stack) == 0 {
		panic("pop from empty stack")
	}
	frame, env.stack = &env.stack[len(env.stack)-1], env.stack[:len(env.stack)-1]
	return frame
}

func (env *Env) setGlobal(key string, value value.Any) {
	env.stack[0].Store(key, value)
}

func (env Env) resolve(sym string) value.Any {
	if len(env.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	// traverse from top of the stack until a binding is found.
	for i := len(env.stack) - 1; i >= 0; i-- {
		if v, found := env.stack[i].Load(sym); found {
			return v
		}
	}

	return nil
}

// Analyzer implementation is responsible for performing syntax analysis
// on given form.
type Analyzer interface {
	// Analyze should perform syntax checks for special forms etc. and
	// return Expr values that can be evaluated against a context.
	Analyze(env *Env, form value.Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(env *Env, form value.Any) (value.Any, error)
}

type stackFrame struct {
	ConcurrentMap

	Name string
	Args []value.Any
}
