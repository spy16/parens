package parens

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spy16/parens/value"
)

var (
	_ = ParseSpecial(parseGoExpr)
	_ = ParseSpecial(parseDefExpr)
	_ = ParseSpecial(parseQuoteExpr)
)

func parseQuoteExpr(_ *Env, args value.Seq) (Expr, error) {
	if count, err := args.Count(); err != nil {
		return nil, err
	} else if count != 1 {
		return nil, Error{
			Cause:   errors.New("invalid quote form"),
			Message: fmt.Sprintf("requires exactly 1 argument, got %d", count),
		}
	}

	first, err := args.First()
	if err != nil {
		return nil, err
	}

	return QuoteExpr{
		Form: first,
	}, nil
}

func parseDefExpr(env *Env, args value.Seq) (Expr, error) {
	if count, err := args.Count(); err != nil {
		return nil, err
	} else if count != 2 {
		return nil, Error{
			Cause:   errors.New("invalid def form"),
			Message: fmt.Sprintf("requires exactly 2 arguments, got %d", count),
		}
	}

	first, err := args.First()
	if err != nil {
		return nil, err
	}

	sym, ok := first.(value.Symbol)
	if !ok {
		return nil, Error{
			Cause:   errors.New("invalid def form"),
			Message: fmt.Sprintf("first arg must be symbol, not '%s'", reflect.TypeOf(first)),
		}
	}

	rest, err := args.Next()
	if err != nil {
		return nil, err
	}

	second, err := rest.First()
	if err != nil {
		return nil, err
	}

	res, err := env.Eval(second)
	if err != nil {
		return nil, err
	}

	return &DefExpr{
		Name:  string(sym),
		Value: res,
	}, nil
}

func parseGoExpr(_ *Env, args value.Seq) (Expr, error) {
	v, err := args.First()
	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, Error{
			Cause: errors.New("go expr requires exactly one argument"),
		}
	}

	return GoExpr{v}, nil
}
