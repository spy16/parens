package parens

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spy16/parens/value"
)

// ErrNotFound is returned by Evaluator when a symbol resolution fails.
var ErrNotFound = errors.New("not found")

// ConcurrentMap to store vars in a stack frame.
type ConcurrentMap interface {
	Load(string) (value.Any, bool)
	Store(string, value.Any)

	// Native map[string]value.Any.  Useful for ranging over values, e.g. when
	// printing a stack frame.
	Native() map[string]value.Any
}

// MapFactory is a constructor for an implementation of ConcurrentMap.
type MapFactory func() ConcurrentMap

// StackFrame .
type StackFrame struct {
	Name string
	Args []value.Any
	ConcurrentMap

	// TODO:  position info?
}

// Context encapsulates execution state.  It is passed to Evaluator in order to evaluate
// a form within a given execution context.
type Context struct {
	stack    []StackFrame
	newMap   MapFactory
	maxDepth int
}

// NewContext returns a global context using the default Context implementation.
// The returned Context IS NOT safe for concurrent use without external synchronization.
func NewContext(opt ...Option) (ctx Context) {
	for _, f := range withDefaults(opt) {
		f(&ctx)
	}

	globalFrame := StackFrame{Name: "<global>", ConcurrentMap: ctx.newMap()}
	ctx.stack = append(ctx.stack, globalFrame)
	return
}

// NewFrame creates an empty StackFrame.  The ConcurrentMap is guaranteed to be non-nil
// and ready for use.  Callers must populate the frame before pushing it to the stack.
func (ctx Context) NewFrame() StackFrame {
	return StackFrame{ConcurrentMap: ctx.newMap()}
}

// Push a frame onto the stack
func (ctx *Context) Push(frame StackFrame) {
	if frame.ConcurrentMap == nil {
		panic("ConcurrentMap must not be nil")
	}

	if len(ctx.stack) == ctx.maxDepth {
		panic(fmt.Sprintf("stack-limit exceeded (maxdepth=%d)", ctx.maxDepth))
	}

	ctx.stack = append(ctx.stack, frame)
}

// Pop a frame from the stack
func (ctx *Context) Pop() (f StackFrame) {
	if len(ctx.stack) == 0 {
		panic("stack must never be empty")
	}

	f, ctx.stack = ctx.stack[len(ctx.stack)-1], ctx.stack[:len(ctx.stack)-1]
	return
}

// SetGlobal variable
func (ctx Context) SetGlobal(name string, val value.Any) {
	ctx.stack[0].Store(name, val)
}

// Resolve name.
func (ctx Context) Resolve(name string) (value.Any, error) {
	if len(ctx.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	// traverse from top of the stack until a binding is found.
	for i := len(ctx.stack) - 1; i >= 0; i-- {
		if v, found := ctx.stack[i].Load(name); found {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

// NewChild derives a copy of the Context
func (ctx Context) NewChild() Context {
	// TODO(performance):  Over-allocate by some reasonable ammount (8-16?).  NewChild is
	// 					   called in preparation for evaluating an InvokeExpr, so we are
	//					   guaranteed to encounter at least one call to ctx.Push().
	//					   InvokeExprs often invoke sub-functions, so we should make a
	//					   reasonable effort to reduce allocations in the common case.
	stack := make([]StackFrame, len(ctx.stack))
	copy(stack, ctx.stack)

	return Context{
		stack:    stack,
		newMap:   ctx.newMap,
		maxDepth: ctx.maxDepth,
	}
}

type basicMap struct {
	sync.RWMutex
	vs map[string]value.Any
}

func newBasicMap() ConcurrentMap {
	return &basicMap{vs: map[string]value.Any{}}
}

func (m *basicMap) Load(name string) (v value.Any, ok bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok = m.vs[name]
	return
}

func (m *basicMap) Store(name string, v value.Any) {
	m.Lock()
	defer m.Unlock()
	m.vs[name] = v
}

func (m *basicMap) Native() map[string]value.Any {
	m.RLock()
	defer m.RUnlock()

	native := make(map[string]value.Any, len(m.vs))
	for k, v := range m.vs {
		native[k] = v
	}

	return native
}
