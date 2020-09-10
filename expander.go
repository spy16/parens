package parens

import "github.com/spy16/parens/value"

type builtinExpander struct{}

func (be *builtinExpander) Expand(ev *Evaluator, form value.Any) (value.Any, error) {
	// TODO: implement macro expansion.
	return nil, nil
}
