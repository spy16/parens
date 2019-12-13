package parens_test

import (
	"github.com/spy16/parens"
	"reflect"
	"testing"
)

func TestInt64_Eval(t *testing.T) {
	testFormEval(t, evalTestCase{
		form:     parens.Int64(10),
		getScope: nil,
		want:     int64(10),
	})
}

func TestFloat64_Eval(t *testing.T) {
	testFormEval(t, evalTestCase{
		form:     parens.Float64(10),
		getScope: nil,
		want:     float64(10),
	})
}

func TestKeyword_Eval(t *testing.T) {
	testFormEval(t, evalTestCase{
		form:     parens.Keyword(":hello"),
		getScope: nil,
		want:     parens.Keyword(":hello"),
	})
}

func TestCharacter_Eval(t *testing.T) {
	testFormEval(t, evalTestCase{
		form:     parens.Character('A'),
		getScope: nil,
		want:     'A',
	})
}

func TestString_Eval(t *testing.T) {
	testFormEval(t, evalTestCase{
		form:     parens.String("hello"),
		getScope: nil,
		want:     "hello",
	})
}

func TestSymbol_Eval(t *testing.T) {
	testAllFormEval(t, []evalTestCase{
		{
			name: "WithBinding",
			form: parens.Symbol("hello"),
			getScope: func() parens.Scope {
				scope := parens.NewScope(nil)
				_ = scope.Bind("hello", parens.Int64(10))
				return scope
			},
			want: parens.Int64(10),
		},
		{
			name: "WithoutBinding",
			form: parens.Symbol("non-existent-symbol"),
			getScope: func() parens.Scope {
				scope := parens.NewScope(nil)
				_ = scope.Bind("hello", parens.Int64(10))
				return scope
			},
			want:    nil,
			wantErr: true,
		},
	})
}

func TestVector_Eval(t *testing.T) {
	testAllFormEval(t, []evalTestCase{
		{
			name: "SimpleVector",
			form: parens.Vector{
				parens.Float64(1.3),
			},
			getScope: nil,
			want: []interface{}{
				float64(1.3),
			},
		},
		{
			name: "VectorWithSymbol",
			form: parens.Vector{
				parens.Float64(1.3),
				parens.Symbol("pi"),
			},
			getScope: func() parens.Scope {
				scope := parens.NewScope(nil)
				_ = scope.Bind("pi", parens.Float64(3.14))
				return scope
			},
			want: []interface{}{
				float64(1.3),
				parens.Float64(3.14),
			},
		},
		{
			name: "VectorWithUnboundSymbol",
			form: parens.Vector{
				parens.Symbol("pi"),
			},
			getScope: func() parens.Scope { return parens.NewScope(nil) },
			want:     []interface{}(nil),
			wantErr:  true,
		},
	})
}

func TestModule_Eval(t *testing.T) {
	testAllFormEval(t, []evalTestCase{
		{
			name: "Literals",
			form: parens.Module{
				parens.String("hello"),
				parens.Int64(10),
				parens.Vector{},
				parens.Keyword("hello"),
				parens.Float64(1.3),
			},
			getScope: nil,
			want:     float64(1.3),
		},
	})
}

func testFormEval(t *testing.T, tt evalTestCase) {
	var scope parens.Scope
	if tt.getScope != nil {
		scope = tt.getScope()
	}
	got, err := tt.form.Eval(scope)
	if (err != nil) != tt.wantErr {
		t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("Eval() got = %#v, want %#v", got, tt.want)
	}
}

func testAllFormEval(t *testing.T, tests []evalTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFormEval(t, tt)
		})
	}
}

type evalTestCase struct {
	name     string
	getScope func() parens.Scope
	form     parens.Expr
	want     interface{}
	wantErr  bool
}
