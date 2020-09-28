package parens_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens"
)

func TestBasicAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		form    parens.Any
		want    parens.Expr
		wantErr bool
	}{
		{
			title: "Nil",
			form:  nil,
			want:  &parens.ConstExpr{Const: parens.Nil{}},
		},
		{
			title: "Symbol",
			form:  parens.Symbol("str"),
			want:  &parens.ConstExpr{Const: parens.String("hello")},
		},
		{
			title:   "Unknown Symbol",
			form:    parens.Symbol("unknown"),
			wantErr: true,
		},
		{
			title: "List With One Entry",
			form:  parens.NewList(parens.Keyword("hello")),
			want: &parens.InvokeExpr{
				Name:   ":hello",
				Target: &parens.ConstExpr{Const: parens.Keyword("hello")},
			},
		},
		{
			title: "List With Multi Entry",
			form:  parens.NewList(parens.Keyword("hello"), parens.Int64(1)),
			want: &parens.InvokeExpr{
				Name:   ":hello",
				Target: &parens.ConstExpr{Const: parens.Keyword("hello")},
				Args: []parens.Expr{
					&parens.ConstExpr{Const: parens.Int64(1)},
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			env := parens.New(parens.WithGlobals(map[string]parens.Any{
				"str": parens.String("hello"),
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
