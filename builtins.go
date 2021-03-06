package parens

import (
	"fmt"
)

var (
	_ Analyzer = (*BuiltinAnalyzer)(nil)
	_ Expander = (*builtinExpander)(nil)
)

// BuiltinAnalyzer parses builtin value forms and returns Expr that can
// be evaluated against parens Env. Custom special form parsers can be
// set using WithAnalyzer().
type BuiltinAnalyzer struct {
	SpecialForms map[string]ParseSpecial
}

// ParseSpecial validates a special form invocation, parse the form and
// returns an expression that can be evaluated for result.
type ParseSpecial func(env *Env, args Seq) (Expr, error)

// Analyze performs syntactic analysis of given form and returns an Expr
// that can be evaluated for result against an Env.
func (ba BuiltinAnalyzer) Analyze(env *Env, form Any) (Expr, error) {
	if IsNil(form) {
		return &ConstExpr{Const: Nil{}}, nil
	}

	switch f := form.(type) {
	case Symbol:
		v := env.Resolve(string(f))
		if v == nil {
			return nil, Error{
				Cause:   ErrNotFound,
				Message: string(f),
			}
		}
		return &ConstExpr{Const: v}, nil

	case Seq:
		cnt, err := f.Count()
		if err != nil {
			return nil, err
		} else if cnt == 0 {
			break
		}

		return ba.analyzeSeq(env, f)
	}

	return &ConstExpr{Const: form}, nil
}

func (ba BuiltinAnalyzer) analyzeSeq(env *Env, seq Seq) (Expr, error) {
	//	Analyze the call target.  This is the first item in the sequence.
	first, err := seq.First()
	if err != nil {
		return nil, err
	}

	// The call target may be a special form.  In this case, we need to get the
	// corresponding parser function, which will take care of parsing/analyzing
	// the tail.
	if sym, ok := first.(Symbol); ok {
		if parse, found := ba.SpecialForms[string(sym)]; found {
			next, err := seq.Next()
			if err != nil {
				return nil, err
			}
			return parse(env, next)
		}
	}

	// Call target is not a special form and must be a Invokable.  Analyze
	// the arguments and create an InvokeExpr.
	ie := InvokeExpr{Name: fmt.Sprintf("%s", first)}
	err = ForEach(seq, func(item Any) (done bool, err error) {
		if ie.Target == nil {
			ie.Target, err = ba.Analyze(env, first)
			return
		}

		var arg Expr
		if arg, err = ba.Analyze(env, item); err == nil {
			ie.Args = append(ie.Args, arg)
		}
		return
	})
	return &ie, err
}

type builtinExpander struct{}

func (be builtinExpander) Expand(_ *Env, _ Any) (Any, error) {
	// TODO: implement macro expansion.
	return nil, nil
}
