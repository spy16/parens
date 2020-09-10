package parens

import (
	"errors"

	"github.com/spy16/parens/value"
)

var _ Context = (*basicContext)(nil)

const globalFrame = "<global>"

// ErrNotFound is returned by Evaluator when a symbol resolution fails.
var ErrNotFound = errors.New("not found")

// Context encapsulates execution state.  It is passed to Evaluator in order to evaluate
// a form within a given execution context.  The default Context implementation IS NOT
// safe for concurrent use without external synchronization.
type Context interface {
	Push(StackFrame)
	Pop() StackFrame

	SetGlobal(name string, val value.Any)
	Resolve(name string) (value.Any, error)

	NewChild() Context
}

// NewContext returns a global context using the default Context implementation.
// The returned Context IS NOT safe for concurrent use without external synchronization.
func NewContext() Context {
	return newContext([]StackFrame{{Name: globalFrame}})
}

// PositionProvider represents the positional information about a value read by reader.
type PositionProvider interface {
	Path() string
	Line() int
	Column() int
	String() string
}

// StackFrame .
type StackFrame struct {
	Name string
	Args []value.Any
	vars map[string]value.Any

	PositionProvider
}

// GetVar from the stack frame
func (frame *StackFrame) GetVar(name string) (value.Any, bool) {
	val, found := frame.vars[name]
	return val, found
}

// SetVar in the stack frame
func (frame *StackFrame) SetVar(name string, val value.Any) {
	if frame.vars == nil {
		frame.vars = map[string]value.Any{}
	}
	frame.vars[name] = val
}

type basicContext struct {
	stack    []StackFrame
	maxDepth int
}

// TODO(enhancement): max stack depth?
func newContext(stack []StackFrame) *basicContext {
	return &basicContext{
		stack: stack,
	}
}

func (ctx *basicContext) Push(frame StackFrame) {
	ctx.stack = append(ctx.stack, frame)
}

func (ctx *basicContext) Pop() (f StackFrame) {
	if len(ctx.stack) == 0 {
		panic("Evaluator stack must never be empty")
	}

	f, ctx.stack = ctx.stack[len(ctx.stack)-1], ctx.stack[:len(ctx.stack)-1]
	return
}

func (ctx basicContext) SetGlobal(name string, val value.Any) {
	ctx.stack[0].SetVar(name, val)
}

func (ctx *basicContext) Resolve(name string) (value.Any, error) {
	if len(ctx.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	// traverse from top of the stack until a binding is found.
	for i := len(ctx.stack) - 1; i >= 0; i-- {
		if v, found := ctx.stack[i].GetVar(name); found {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

func (ctx basicContext) NewChild() Context {
	stack := make([]StackFrame, len(ctx.stack))
	copy(stack, ctx.stack)
	return &basicContext{stack: stack}
}
