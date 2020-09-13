package parens

import (
	"sync"

	"github.com/spy16/parens/value"
)

var _ ConcurrentMap = (*mutexMap)(nil)

// ConcurrentMap implementation can be set on the root context to customise
// the map used for storing variables in the stack frames.
type ConcurrentMap interface {
	// Store should store the key-value pair in the map.
	Store(key string, val value.Any)

	// Load should return the value associated with the key if it exists.
	// Returns nil, false otherwise.
	Load(key string) (value.Any, bool)

	// Map should return a native Go map of all key-values in the concurrent
	// map. This can be used for iteration etc.
	Map() map[string]value.Any
}

// mutexMap implements a simple ConcurrentMap using sync.RWMutex locks. Zero
// value is ready for use.
type mutexMap struct {
	sync.RWMutex
	vs map[string]value.Any
}

func (m *mutexMap) Load(name string) (v value.Any, ok bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok = m.vs[name]
	return
}

func (m *mutexMap) Store(name string, v value.Any) {
	m.Lock()
	defer m.Unlock()

	if m.vs == nil {
		m.vs = map[string]value.Any{}
	}
	m.vs[name] = v
}

func (m *mutexMap) Map() map[string]value.Any {
	m.RLock()
	defer m.RUnlock()

	native := make(map[string]value.Any, len(m.vs))
	for k, v := range m.vs {
		native[k] = v
	}

	return native
}
