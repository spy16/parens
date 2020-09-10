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
	first, err := seq.First()
	if err != nil {
		return nil, err
	}

	// handle special form analysis.
	if sym, ok := first.(*value.Symbol); ok {
		parse, found := ba.specials[sym.Value]
		if found {
			next, err := seq.Next()
			if err != nil {
				return nil, err
			}
			return parse(ev, next)
		}
	}

	target, err := ba.Analyze(ev, first)
	if err != nil {
		return nil, err
	}

	var args []Expr
	for seq != nil {
		seq, err = seq.Next()
		if err != nil {
			return nil, err
		}

		f, err := seq.First()
		if err != nil {
			return nil, err
		}

		arg, err := ba.Analyze(ev, f)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return InvokeExpr{
		Name:   fmt.Sprintf("%s", target),
		Target: target,
		Args:   args,
	}, nil
}
