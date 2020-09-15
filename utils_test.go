package parens_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens"
	"github.com/spy16/parens/value"
)

func TestEvalAll(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		env     *parens.Env
		vals    []value.Any
		want    []value.Any
		wantErr bool
	}{
		{
			title: "EmptyList",
			env:   parens.New(),
			vals:  nil,
			want:  []value.Any{},
		},
		{
			title:   "EvalFails",
			env:     parens.New(),
			vals:    []value.Any{value.Symbol("foo")},
			wantErr: true,
		},
		{
			title: "Success",
			env:   parens.New(),
			vals:  []value.Any{value.String("foo"), value.Keyword("hello")},
			want:  []value.Any{value.String("foo"), value.Keyword("hello")},
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			got, err := parens.EvalAll(tt.env, tt.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalAll() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvalAll() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
