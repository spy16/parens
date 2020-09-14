package parens

import "github.com/spy16/parens/value"

type basicExpander struct{}

func (be basicExpander) Expand(_ *Env, _ value.Any) (value.Any, error) {
	// TODO: implement macro expansion.
	return nil, nil
}
