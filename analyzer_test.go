package parens_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens"
	"github.com/spy16/parens/value"
)

func TestBasicAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		form    value.Any
		want    parens.Expr
		wantErr bool
	}{
		{
			title: "Nil",
			form:  nil,
			want:  &parens.ConstExpr{Const: value.Nil{}},
		},
		{
			title: "Symbol",
			form:  value.Symbol("str"),
			want:  &parens.ConstExpr{Const: value.String("hello")},
		},
		{
			title:   "Unknown Symbol",
			form:    value.Symbol("unknown"),
			wantErr: true,
		},
		{
			title: "List",
			form:  value.NewList(value.String("hello")),
			want: &parens.InvokeExpr{
				Name:   "hello",
				Target: &parens.ConstExpr{Const: value.String("hello")},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			env := parens.New(parens.WithGlobals(map[string]value.Any{
				"str": value.String("hello"),
			}, nil))

			az := &parens.BuiltinAnalyzer{}
			got, err := az.Analyze(env, tt.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuiltinAnalyzer.Analyze() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuiltinAnalyzer.Analyze() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
