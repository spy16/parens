package parens

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func add(a, b float64) float64 {
	return a + b
}

func BenchmarkParens_Execute(suite *testing.B) {
	env := NewScope(nil)
	suite.Run("Execute", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ExecuteStr("(add 1 2)", env)
		}
	})

	expr := List(
		[]Expr{
			Symbol("add"),
			Int64(1),
			Int64(2),
		})

	suite.Run("ExecuteExpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ExecuteExpr(expr, env)
		}
	})
}

func BenchmarkParens_FunctionCall(suite *testing.B) {
	suite.Run("DirectCall", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(1, 2)
		}
	})

	expr, err := Parse(strings.NewReader("(add 1 2)"))
	if err != nil {
		suite.Fatalf("failed to parse expression: %s", err)
	}

	suite.Run("CallThroughParens", func(b *testing.B) {
		scope := NewScope(nil)
		scope.Bind("add", add)

		for i := 0; i < b.N; i++ {
			expr.Eval(scope)
		}
	})
}

func TestExecute_Success(t *testing.T) {
	scope := NewScope(nil)

	res, err := Execute(strings.NewReader("10"), scope)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, int64(10), res)
}

func TestExecute_EvalFailure(t *testing.T) {
	scope := NewScope(nil)

	res, err := Execute(strings.NewReader("(hello)"), scope)
	require.Error(t, err)
	assert.Nil(t, res)
}
