package parens_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/spy16/parens"
	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func add(a, b float64) float64 {
	return a + b
}

func BenchmarkParens_Execute(suite *testing.B) {
	ins := parens.New(parens.NewScope(nil))
	suite.Run("Execute", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ins.Execute("(add 1 2)")
		}
	})

	expr := parser.ListExpr{
		List: []parser.Expr{
			parser.SymbolExpr{
				Symbol: "add",
			},
			parser.NumberExpr{
				Number: 1,
			},
			parser.NumberExpr{
				Number: 2,
			},
		},
	}

	suite.Run("ExecuteExpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ins.ExecuteExpr(expr)
		}
	})
}

func BenchmarkParens_FunctionCall(suite *testing.B) {
	suite.Run("DirectCall", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(1, 2)
		}
	})

	expr, err := parser.Parse(strings.NewReader("(add 1 2)"))
	if err != nil {
		suite.Fatalf("failed to parse expression: %s", err)
	}

	suite.Run("CallThroughParens", func(b *testing.B) {
		scope := parens.NewScope(nil)
		scope.Bind("add", add)

		for i := 0; i < b.N; i++ {
			expr.Eval(scope)
		}
	})
}

func TestExecute_Success(t *testing.T) {
	scope := parens.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockExpr(10, nil), nil)

	res, err := par.Execute("10")
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res, 10)
}

func TestExecute_EvalFailure(t *testing.T) {
	scope := parens.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockExpr(nil, errors.New("failed")), nil)

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func TestExecute_ParseFailure(t *testing.T) {
	scope := parens.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(nil, errors.New("failed"))

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func mockExpr(v interface{}, err error) parser.Expr {
	return exprMock(func(scope parser.Scope) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}

func mockParseFn(expr parser.Expr, err error) parens.ParseFn {
	return func(name string, rd io.RuneScanner) (parser.Expr, error) {
		if err != nil {
			return nil, err
		}
		return expr, nil
	}
}

type exprMock func(scope parser.Scope) (interface{}, error)

func (sm exprMock) Eval(scope parser.Scope) (interface{}, error) {
	return sm(scope)
}
