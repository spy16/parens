package parens

import (
	"fmt"
	"strings"
)

var _ Expander = (*builtinExpander)(nil)

// Macro is recognized by the built-in expander and expanded into an expression.
type Macro struct {
	Name string   // macro name.
	Args []string // argument names.
	Body Any
}

// SExpr returns a valid s-expression for Macro.
func (m Macro) SExpr() (string, error) {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(m.Name)
	for _, arg := range m.Args {
		b.WriteString(arg)
	}
	b.WriteRune(')')
	return b.String(), nil
}

type builtinExpander struct{}

func (be builtinExpander) Expand(env *Env, form Any) (Any, error) {
	if seq, ok := form.(Seq); ok {
		first, err := seq.First()
		if err != nil {
			return nil, err
		}

		switch f := first.(type) {
		case Macro:
			return be.expand(env, f, seq)

		case Symbol:
			if m, ok := env.Resolve(f.String()).(Macro); ok {
				if seq, err = seq.Next(); err != nil {
					return nil, err
				}

				return be.expand(env, m, seq)
			}
		}
	}

	return nil, nil
}

func (be builtinExpander) expand(env *Env, m Macro, seq Seq) (Any, error) {
	n, err := seq.Count()
	if err != nil {
		return nil, err
	}

	if n != len(m.Args) {
		return nil, Error{
			Cause: fmt.Errorf("expected %d arguments, got %d", len(m.Args), n),
		}
	}

	// TODO(enhancement):  handle symbol collision

	var f stackFrame
	f.Name = m.Name
	// f.Args = make([]Any, n)  // is this needed?
	f.Vars = make(map[string]Any, n)

	for _, sym := range m.Args {
		arg, err := seq.First()
		if err != nil {
			return nil, err
		}

		f.Vars[sym] = arg

		if seq, err = seq.Next(); err != nil {
			return nil, err
		}
	}

	env.push(f) // Where does this get popped?

	return env.Eval(m.Body)

	// 1. Create bindings for parameters and the corresponding invocation arguments
	//	  in the Env. (symbol collision needs to be handled. See gensym).
	// 2. Then evaluate the m.Body in current Env.
	// 3. Return the result as the new form.
}
