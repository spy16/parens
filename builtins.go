package parens

// ConstExpr returns the Const value wrapped inside when evaluated. It has
// no side-effect on the VM.
type ConstExpr struct{ Const Any }

func (ce ConstExpr) Eval() (Any, error) { return ce.Const, nil }

type builtinAnalyzer struct {
	extender Analyzer
}

func (ba *builtinAnalyzer) Analyze(form Any) (Expr, error) {
	if ba.extender != nil {
		return ba.extender.Analyze(form)
	}

	return &ConstExpr{Const: form}, nil
}

type builtinExpander struct{}

func (be *builtinExpander) Expand(p *Evaluator, form Any) (Any, error) {
	// TODO: implement macro expansion.
	return nil, nil
}
