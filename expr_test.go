package parens_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reader"
	"github.com/spy16/parens/value"
)

func TestGoExpr(t *testing.T) {
	env := parens.New()

	b := bytes.NewBufferString("(go (def test :keyword))")
	r := reader.New(b)

	val, err := r.One()
	if err != nil {
		t.Error(err)
		return
	}

	env.Eval(val)
	time.Sleep(time.Millisecond)

	val, err = env.Eval(value.Symbol("test"))
	if err != nil {
		t.Error(err)
		return
	}

	kw, ok := val.(value.Keyword)
	if !ok {
		t.Errorf("expected value.Keyword, got %s", reflect.TypeOf(kw))
		return
	}

	if string(kw) != "keyword" {
		t.Errorf("expected keyword value of \"keyword\", got \"%s\"", string(kw))
		return
	}
}
