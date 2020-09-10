package parens

import (
	"errors"

	"github.com/spy16/parens/value"
)

var (
	_ = ParseSpecial(ParseGo)
)

// ParseGo parses a special form into a GoExpr.
func ParseGo(ev Evaluator, args value.Seq) (Expr, error) {
	v, err := args.First()
	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, NewSyntaxError(errors.New("go expr requires exactly one argument"))
	}

	return GoExpr{v}, nil
}
