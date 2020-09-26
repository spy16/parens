package parens_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reader"
)

func TestDoExpr_Eval(t *testing.T) {
	t.Parallel()

	t.Run("No Body", func(t *testing.T) {
		de := parens.DoExpr{}
		res, err := de.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.Nil{}, res)
	})

	t.Run("With Body", func(t *testing.T) {
		de := parens.DoExpr{
			Exprs: []parens.Expr{
				&parens.ConstExpr{Const: parens.Symbol("foo")},
			},
		}
		res, err := de.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.Symbol("foo"), res)
	})
}

func TestIfExpr_Eval(t *testing.T) {
	t.Parallel()

	t.Run("If Nil", func(t *testing.T) {
		ie := parens.IfExpr{
			Test: nil,
			Then: &parens.ConstExpr{Const: parens.String("then")},
			Else: &parens.ConstExpr{Const: parens.String("else")},
		}
		res, err := ie.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.String("else"), res)
	})

	t.Run("If True", func(t *testing.T) {
		ie := parens.IfExpr{
			Test: &parens.ConstExpr{Const: parens.Bool(true)},
			Then: &parens.ConstExpr{Const: parens.String("then")},
			Else: &parens.ConstExpr{Const: parens.String("else")},
		}
		res, err := ie.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.String("then"), res)
	})

	t.Run("If Val", func(t *testing.T) {
		ie := parens.IfExpr{
			Test: &parens.ConstExpr{Const: parens.String("foo")},
			Then: &parens.ConstExpr{Const: parens.String("then")},
			Else: &parens.ConstExpr{Const: parens.String("else")},
		}
		res, err := ie.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.String("then"), res)
	})

	t.Run("No Body", func(t *testing.T) {
		ie := parens.IfExpr{
			Test: &parens.ConstExpr{Const: parens.String("foo")},
		}
		res, err := ie.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.Nil{}, res)
	})
}

func TestDefExpr_Eval(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Name", func(t *testing.T) {
		de := parens.DefExpr{
			Env:   parens.New(),
			Name:  "",
			Value: &parens.ConstExpr{Const: parens.Int64(10)},
		}
		v, err := de.Eval()
		assertErr(t, err)
		assertEqual(t, nil, v)
	})

	t.Run("Success", func(t *testing.T) {
		de := parens.DefExpr{
			Env:   parens.New(),
			Name:  "foo",
			Value: &parens.ConstExpr{Const: parens.Int64(10)},
		}
		v, err := de.Eval()
		requireNoErr(t, err)
		assertEqual(t, parens.Symbol("foo"), v)
	})
}

func TestQuoteExpr_Eval(t *testing.T) {
	want := parens.NewList()

	qe := parens.QuoteExpr{Form: want}
	got, err := qe.Eval()
	requireNoErr(t, err)

	assertEqual(t, want, got)
}

func TestGoExpr_Eval(t *testing.T) {
	r := reader.New(strings.NewReader("(go (def test :keyword))"))
	actual, err := r.One()
	requireNoErr(t, err)

	env := parens.New()
	_, _ = env.Eval(actual)
	time.Sleep(5 * time.Millisecond)

	actual, err = env.Eval(parens.Symbol("test"))
	requireNoErr(t, err)

	if kw, ok := actual.(parens.Keyword); !ok {
		t.Errorf("expected parens.Keyword, got %s", reflect.TypeOf(kw))
		return
	} else if string(kw) != "keyword" {
		t.Errorf("expected keyword value of \"keyword\", got \"%s\"", string(kw))
		return
	}
}

func requireNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertEqual(t *testing.T, want interface{}, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got=%#v\nwant=%#v", got, want)
	}
}

func assertErr(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expecting error, got nil")
	}
}
