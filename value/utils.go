package value

import (
	"reflect"
	"strings"
)

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
		sexpr, err := item.SExpr()
		if err != nil {
			return false, err
		}

		b.WriteString(sexpr)
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
