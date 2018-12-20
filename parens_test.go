package parens_test

import (
	"strings"
	"testing"

	"github.com/spy16/parens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func add(a, b float64) float64 {
	return a + b
}

func BenchmarkParens_Execute(suite *testing.B) {
	env := parens.NewScope(nil)
	suite.Run("Execute", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parens.ExecuteStr("(add 1 2)", env)
		}
	})

	expr := parens.ListExpr(
		[]parens.Expr{
			parens.SymbolExpr("add"),
			parens.NumberExpr{
				Number: 1,
			},
			parens.NumberExpr{
				Number: 2,
			},
		})

	suite.Run("ExecuteExpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parens.ExecuteExpr(expr, env)
		}
	})
}

func BenchmarkParens_FunctionCall(suite *testing.B) {
	suite.Run("DirectCall", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(1, 2)
		}
	})

	expr, err := parens.Parse(strings.NewReader("(add 1 2)"))
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

	res, err := parens.ExecuteOne(strings.NewReader("10"), scope)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 10.0, res)
}

func TestExecute_EvalFailure(t *testing.T) {
	scope := parens.NewScope(nil)

	res, err := parens.ExecuteOne(strings.NewReader("(hello)"), scope)
	require.Error(t, err)
	assert.Nil(t, res)
}
