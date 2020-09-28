package parens

var _ Expander = (*builtinExpander)(nil)

// Macro is recognized by the built-in expander and expanded into an expression.
type Macro struct {
	Name string   // macro name.
	Args []string // argument names.
	Body Any
}

type builtinExpander struct{}

func (be builtinExpander) Expand(_ *Env, _ Any) (Any, error) {
	// TODO: implement macro expansion.
	return nil, nil
}
