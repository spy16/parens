package parens

import (
	"fmt"

	"github.com/spy16/parens/value"
)

// ParseSpecial validates a special form invocation, parse the form and
// returns an expression that can be evaluated for result.
type ParseSpecial func(ev *Evaluator, args value.Seq) (Expr, error)

type builtinAnalyzer struct {
	ev       *Evaluator
	extender Analyzer
	specials map[string]ParseSpecial
}

func (ba *builtinAnalyzer) Analyze(form value.Any) (Expr, error) {
	if seq, isSeq := form.(value.Seq); isSeq {
		return ba.analyzeSeq(seq)
	}

	if ba.extender != nil {
		return ba.extender.Analyze(form)
	}

	return &ConstExpr{Const: form}, nil
}

func (ba *builtinAnalyzer) analyzeSeq(seq value.Seq) (Expr, error) {
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
			return parse(ba.ev, next)
		}
	}

	target, err := ba.Analyze(first)
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

		arg, err := ba.Analyze(f)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return &InvokeExpr{
		Name:   fmt.Sprintf("%s", target),
		Target: target,
		Args:   args,
	}, nil
}
