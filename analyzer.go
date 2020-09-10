package parens

import (
	"fmt"

	"github.com/spy16/parens/value"
)

var _ Analyzer = (*BasicAnalyzer)(nil)

// ParseSpecial validates a special form invocation, parse the form and
// returns an expression that can be evaluated for result.
type ParseSpecial func(ev Evaluator, args value.Seq) (Expr, error)

// BasicAnalyzer can parse (optional) special forms.
type BasicAnalyzer struct {
	specials map[string]ParseSpecial
}

// Analyze the form.
func (ba BasicAnalyzer) Analyze(ev Evaluator, form value.Any) (Expr, error) {
	switch f := form.(type) {
	case value.Seq:
		cnt, err := f.Count()
		if err != nil {
			return nil, err
		} else if cnt == 0 {
			break
		}

		return ba.analyzeSeq(ev, f)
	}

	return ConstExpr{Const: form}, nil
}

func (ba BasicAnalyzer) analyzeSeq(ev Evaluator, seq value.Seq) (Expr, error) {
	/*
		First we analyze the call target.  This is the first item in the sequence.
	*/
	first, err := seq.First()
	if err != nil {
		return nil, err
	}

	/*
		The call target may be a special form.  In this case, we need to get the
		corresponding parser function, which will take care of parsing/analyzing the
		tail.
	*/
	if sym, ok := first.(*value.Symbol); ok {
		if parse, found := ba.specials[sym.Value]; found {
			next, err := seq.Next()
			if err != nil {
				return nil, err
			}

			return parse(ev, next)
		}
	}

	/*
		If we get here, the call target is a standard invokable (usually a function or
		a macro), so we are responsible for analyzing the call target and its arguments.
	*/
	var target Expr
	var args []Expr
	err = value.ForEach(seq, func(item value.Any) (done bool, err error) {
		if target == nil {
			if target, err = ba.Analyze(ev, first); err != nil {
				return
			}
		}

		var arg Expr
		if arg, err = ba.Analyze(ev, item); err == nil {
			args = append(args, arg)
		}

		return
	})

	return InvokeExpr{
		Name:   fmt.Sprintf("%s", target),
		Target: target,
		Args:   args,
	}, err
}
