package parens

import (
	"fmt"
	"reflect"
	"strings"
)

// EvalAll evaluates each value in the list against the given env and returns
// a list of resultant values.
func EvalAll(env *Env, vals []Any) ([]Any, error) {
	res := make([]Any, 0, len(vals))
	for _, form := range vals {
		form, err := env.Eval(form)
		if err != nil {
			return nil, err
		}
		res = append(res, form)
	}
	return res, nil
}

// IsNil returns true if value is native go `nil` or `Nil{}`.
func IsNil(v Any) bool {
	if v == nil {
		return true
	}
	_, isNilType := v.(Nil)
	return isNilType
}

// IsTruthy returns true if the value has a logical vale of `true`.
func IsTruthy(v Any) bool {
	if IsNil(v) {
		return false
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Bool {
		return rv.Bool()
	}

	return true
}

// SeqString returns a string representation for the sequence with given prefix
// suffix and separator.
func SeqString(seq Seq, begin, end, sep string) (string, error) {
	var b strings.Builder
	b.WriteString(begin)
	err := ForEach(seq, func(item Any) (bool, error) {
		if sxpr, ok := item.(SExpressable); ok {
			s, err := sxpr.SExpr()
			if err != nil {
				return false, err
			}
			b.WriteString(s)

		} else {
			b.WriteString(fmt.Sprintf("%#v", item))
		}

		b.WriteString(sep)
		return false, nil
	})

	if err != nil {
		return "", err
	}

	return strings.TrimRight(b.String(), sep) + end, err
}

// ForEach reads from the sequence and calls the given function for each item.
// Function can return true to stop the iteration.
func ForEach(seq Seq, call func(item Any) (bool, error)) (err error) {
	var v Any
	var done bool
	for seq != nil {
		if v, err = seq.First(); err != nil || v == nil {
			break
		}

		if done, err = call(v); err != nil || done {
			break
		}

		if seq, err = seq.Next(); err != nil {
			break
		}
	}

	return
}
