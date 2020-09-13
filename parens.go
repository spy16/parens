package parens

import (
	"context"

	"github.com/spy16/parens/value"
)

// New returns a new root context initialised based on given options.
func New(opts ...Option) *Context {
	p := &Context{}
	for _, opt := range withDefaults(opts) {
		opt(p)
	}
	return p
}

// EvalAll evaluates each value in the list against the given env and returns a list
// of resultant value.
func EvalAll(ctx *Context, vals []value.Any) (res []value.Any, err error) {
	res = make([]value.Any, 0, len(vals))

	for _, form := range vals {
		if form, err = ctx.Eval(form); err != nil {
			break
		}

		res = append(res, form)
	}

	return
}

// Context represents the environment/context in which forms are evaluated
// for result. Context is not safe for concurrent use. Use fork() to get a
// child context for concurrent executions.
type Context struct {
	ctx        context.Context
	parent     *Context
	analyzer   Analyzer
	expander   Expander
	stack      []stackFrame
	maxDepth   int
	mapFactory func() ConcurrentMap
}

// Eval performs macro-expansion if necessary, converts the expanded form
// to an expression and evaluates the resulting expression.
func (ctx *Context) Eval(form value.Any) (value.Any, error) {
	expr, err := ctx.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return nil, nil
	}

	return expr.Eval(ctx)
}

func (ctx *Context) expandAnalyze(form value.Any) (Expr, error) {
	if expr, ok := form.(Expr); ok {
		// Already an Expr, nothing to do.
		return expr, nil
	}

	if expanded, err := ctx.expander.Expand(ctx, form); err != nil {
		return nil, err
	} else if expanded != nil {
		// Expansion did happen. Throw away the old form and continue with
		// the expanded version.
		form = expanded
	}

	return ctx.analyzer.Analyze(ctx, form)
}

// fork creates a child context from the parent and returns which can be
// used as context for an independent thread of execution.
func (ctx *Context) fork() *Context {
	return &Context{
		ctx:        ctx.ctx,
		parent:     ctx,
		stack:      append([]stackFrame(nil), ctx.stack...),
		expander:   ctx.expander,
		analyzer:   ctx.analyzer,
		maxDepth:   ctx.maxDepth,
		mapFactory: ctx.mapFactory,
	}
}

func (ctx *Context) push(frame stackFrame) {
	ctx.stack = append(ctx.stack, frame)
}

func (ctx *Context) pop() (frame *stackFrame) {
	if len(ctx.stack) == 0 {
		panic("pop from empty stack")
	}
	frame, ctx.stack = &ctx.stack[len(ctx.stack)-1], ctx.stack[1:]
	return frame
}

func (ctx *Context) setGlobal(key string, value value.Any) {
	rootCtx := ctx
	for rootCtx.parent != nil {
		rootCtx = ctx.parent
	}

	// TODO: verify this. what is expected if this is a child context?
	rootCtx.stack[0].Store(key, value)
}

// Analyzer implementation is responsible for performing syntax analysis
// on given form.
type Analyzer interface {
	// Analyze should perform syntax checks for special forms etc. and
	// return Expr values that can be evaluated against a context.
	Analyze(ctx *Context, form value.Any) (Expr, error)
}

// Expander implementation is responsible for performing macro-expansion
// where necessary.
type Expander interface {
	// Expand should expand/rewrite the given form if it's a macro form
	// and return the expanded version. If given form is not macro form,
	// it should return nil, nil.
	Expand(p *Context, form value.Any) (value.Any, error)
}

type stackFrame struct {
	Name string
	Args []value.Any
	ConcurrentMap
}
