package parens

import (
	"context"
	"sync"
)

var _ ConcurrentMap = (*mutexMap)(nil)

// Env represents the environment/context in which forms are evaluated
// for result. Env is not safe for concurrent use. Use fork() to get a
// child context for concurrent executions.
type Env struct {
	ctx      context.Context
	analyzer Analyzer
	expander Expander
	globals  ConcurrentMap
	stack    []stackFrame
	maxDepth int
}

// ConcurrentMap is used by the Env to store variables in the global stack frame.
type ConcurrentMap interface {
	// Store should store the key-value pair in the map.
	Store(key string, val Any)

	// Load should return the value associated with the key if it exists.
	// Returns nil, false otherwise.
	Load(key string) (Any, bool)

	// Map should return a native Go map of all key-values in the concurrent
	// map. This can be used for iteration etc.
	Map() map[string]Any
}

// Eval performs macro-expansion if necessary, converts the expanded form
// to an expression and evaluates the resulting expression.
func (env *Env) Eval(form Any) (Any, error) {
	expr, err := env.expandAnalyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return nil, nil
	}

	return expr.Eval(env)
}

// Resolve a symbol.
func (env Env) Resolve(sym string) Any {
	if len(env.stack) > 0 {
		// check inside top of the stack for local bindings.
		top := env.stack[len(env.stack)-1]
		if v, found := top.Vars[sym]; found {
			return v
		}
	}
	// return the value from global bindings if found.
	v, _ := env.globals.Load(sym)
	return v
}

func (env *Env) expandAnalyze(form Any) (Expr, error) {
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

// Fork creates a child context from Env and returns it. The child context
// can be used as context for an independent thread of execution.
func (env *Env) Fork() *Env {
	return &Env{
		ctx:      env.ctx,
		globals:  env.globals,
		expander: env.expander,
		analyzer: env.analyzer,
		maxDepth: env.maxDepth,
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

func (env *Env) setGlobal(key string, value Any) {
	env.globals.Store(key, value)
}

type stackFrame struct {
	Name string
	Args []Any
	Vars map[string]Any
}

func newMutexMap() ConcurrentMap { return &mutexMap{} }

// mutexMap implements a simple ConcurrentMap using sync.RWMutex locks. Zero
// value is ready for use.
type mutexMap struct {
	sync.RWMutex
	vs map[string]Any
}

func (m *mutexMap) Load(name string) (v Any, ok bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok = m.vs[name]
	return
}

func (m *mutexMap) Store(name string, v Any) {
	m.Lock()
	defer m.Unlock()

	if m.vs == nil {
		m.vs = map[string]Any{}
	}
	m.vs[name] = v
}

func (m *mutexMap) Map() map[string]Any {
	m.RLock()
	defer m.RUnlock()

	native := make(map[string]Any, len(m.vs))
	for k, v := range m.vs {
		native[k] = v
	}

	return native
}
