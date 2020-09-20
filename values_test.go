package parens_test

import (
	"testing"

	"github.com/spy16/parens"
)

func TestComp(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		a, b    parens.Any
		op      func(a, b parens.Any) (bool, error)
		want    bool
		wantErr bool
	}{
		// Nil
		{
			desc: "nil == nil",
			a:    parens.Nil{},
			b:    parens.Nil{},
			op:   parens.Eq,
			want: true,
		},
		{
			desc: "nil == false",
			a:    parens.Nil{},
			b:    parens.Bool(false),
			op:   parens.Eq,
		},

		// Bool
		{
			desc: "true == true",
			a:    parens.Bool(true),
			b:    parens.Bool(true),
			op:   parens.Eq,
			want: true,
		},
		{
			desc: "false == false",
			a:    parens.Bool(true),
			b:    parens.Bool(true),
			op:   parens.Eq,
			want: true,
		},
		{
			desc: "false == true",
			a:    parens.Bool(false),
			b:    parens.Bool(true),
			op:   parens.Eq,
		},

		// Int64
		{
			desc: "0 == 0",
			a:    parens.Int64(0),
			b:    parens.Int64(0),
			op:   parens.Eq,
			want: true,
		},
		{
			desc: "0 < 1",
			a:    parens.Int64(0),
			b:    parens.Int64(1),
			op:   parens.Lt,
			want: true,
		},
		{
			desc: "1 > 0",
			a:    parens.Int64(1),
			b:    parens.Int64(0),
			op:   parens.Gt,
			want: true,
		},
		{
			desc: "0 <= 1",
			a:    parens.Int64(0),
			b:    parens.Int64(1),
			op:   parens.Le,
			want: true,
		},
		{
			desc: "1 >= 0",
			a:    parens.Int64(1),
			b:    parens.Int64(0),
			op:   parens.Ge,
			want: true,
		},
		{
			desc: "0 <= 0",
			a:    parens.Int64(0),
			b:    parens.Int64(0),
			op:   parens.Le,
			want: true,
		},
		{
			desc: "0 >= 0",
			a:    parens.Int64(0),
			b:    parens.Int64(0),
			op:   parens.Ge,
			want: true,
		},

		// Float64
		{
			desc: "0. == 0.",
			a:    parens.Float64(0),
			b:    parens.Float64(0),
			op:   parens.Eq,
			want: true,
		},
		{
			desc: "0. < 1.",
			a:    parens.Float64(0),
			b:    parens.Float64(1),
			op:   parens.Lt,
			want: true,
		},
		{
			desc: "1. > 0.",
			a:    parens.Float64(1),
			b:    parens.Float64(0),
			op:   parens.Gt,
			want: true,
		},
		{
			desc: "0. <= 1.",
			a:    parens.Float64(0),
			b:    parens.Float64(1),
			op:   parens.Le,
			want: true,
		},
		{
			desc: "1. >= 0.",
			a:    parens.Float64(1),
			b:    parens.Float64(0),
			op:   parens.Ge,
			want: true,
		},
		{
			desc: "0. <= 0.",
			a:    parens.Float64(0),
			b:    parens.Float64(0),
			op:   parens.Le,
			want: true,
		},
		{
			desc: "0. >= 0.",
			a:    parens.Float64(0),
			b:    parens.Float64(0),
			op:   parens.Ge,
			want: true,
		},
		{
			// LinkedList
			desc: "(1 2 3) == (1 2 3)",
			a:    parens.NewList(parens.Int64(1), parens.Int64(2), parens.Int64(3)),
			b:    parens.NewList(parens.Int64(1), parens.Int64(2), parens.Int64(3)),
			op:   parens.Eq,
			want: true,
		},
		{
			// LinkedList
			desc: "(1 2 3) == (1 2 nil)",
			a:    parens.NewList(parens.Int64(1), parens.Int64(2), parens.Int64(3)),
			b:    parens.NewList(parens.Int64(1), parens.Int64(2), parens.Nil{}),
			op:   parens.Eq,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := tt.op(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %#v, wantErr %#v", tt.desc, err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("%s returned %t, expected %t", tt.desc, got, tt.want)
			}
		})
	}
}
