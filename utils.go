package parens

import "github.com/spy16/parens/value"

// EvalAll evaluates each value in the list against the given env and returns
// a list of resultant values.
func EvalAll(env *Env, vals []value.Any) ([]value.Any, error) {
	res := make([]value.Any, 0, len(vals))
	for _, form := range vals {
		form, err := env.Eval(form)
		if err != nil {
			return nil, err
		}
		res = append(res, form)
	}
	return res, nil
}
