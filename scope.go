package parens

import (
	"fmt"
	"strings"

	"github.com/spy16/parens/reflection"
)

// NewScope initializes a new scope with given parent scope. parent
// can be nil.
func NewScope(parent Scope) Scope {
	return &defaultScope{
		parent: parent,
		vals:   map[string]scopeEntry{},
	}
}

type defaultScope struct {
	parent Scope
	vals   map[string]scopeEntry
}

type scopeEntry struct {
	val reflection.Value
	doc string
}

func (sc *defaultScope) Root() Scope {
	if sc.parent == nil {
		return sc
	}

	return sc.parent.Root()
}

func (sc *defaultScope) Bind(name string, v interface{}, doc ...string) error {
	val := reflection.NewValue(v)
	sc.vals[name] = scopeEntry{
		val: val,
		doc: strings.TrimSpace(strings.Join(doc, "\n")),
	}

	return nil
}

func (sc *defaultScope) Doc(name string) string {
	if entry := sc.entry(name); entry != nil {
		return entry.doc
	}

	if sc.parent != nil {
		if swd, ok := sc.parent.(scopeWithDoc); ok {
			return swd.Doc(name)
		}
	}

	return ""
}

func (sc *defaultScope) Get(name string) (interface{}, error) {
	entry := sc.entry(name)
	if entry == nil {
		if sc.parent != nil {
			return sc.parent.Get(name)
		}
		return nil, fmt.Errorf("name '%s' not found", name)
	}

	return entry.val.RVal.Interface(), nil
}

func (sc *defaultScope) String() string {
	str := []string{}
	for name := range sc.vals {
		str = append(str, fmt.Sprintf("%s", name))
	}
	return strings.Join(str, "\n")
}

func (sc *defaultScope) entry(name string) *scopeEntry {
	entry, found := sc.vals[name]
	if found {
		return &entry
	}

	return nil
}

type scopeWithDoc interface {
	Doc(name string) string
}
