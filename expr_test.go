package parens_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reader"
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

	_, _ = env.Eval(val)
	time.Sleep(time.Millisecond)

	val, err = env.Eval(parens.Symbol("test"))
	if err != nil {
		t.Error(err)
		return
	}

	kw, ok := val.(parens.Keyword)
	if !ok {
		t.Errorf("expected parens.Keyword, got %s", reflect.TypeOf(kw))
		return
	}

	if string(kw) != "keyword" {
		t.Errorf("expected keyword value of \"keyword\", got \"%s\"", string(kw))
		return
	}
}
