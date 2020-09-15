package parens_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens"
)

func TestEvalAll(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		env     *parens.Env
		vals    []parens.Any
		want    []parens.Any
		wantErr bool
	}{
		{
			title: "EmptyList",
			env:   parens.New(),
			vals:  nil,
			want:  []parens.Any{},
		},
		{
			title:   "EvalFails",
			env:     parens.New(),
			vals:    []parens.Any{parens.Symbol("foo")},
			wantErr: true,
		},
		{
			title: "Success",
			env:   parens.New(),
			vals:  []parens.Any{parens.String("foo"), parens.Keyword("hello")},
			want:  []parens.Any{parens.String("foo"), parens.Keyword("hello")},
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
