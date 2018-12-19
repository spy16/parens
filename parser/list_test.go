package parser_test

import (
	"io"
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_List(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parser.Parse(reader(")"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parser.Parse(reader("("))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("EmptyList", func(t *testing.T) {
		expr, err := parser.Parse(reader("()"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.ListExpr{}, expr)
		assert.Equal(t, 0, len(expr.(parser.ListExpr).List))
	})

	suite.Run("SimpleList", func(t *testing.T) {
		expr, err := parser.Parse(reader("(1 \"hello\")"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parser.ListExpr).List))
	})

	suite.Run("NestedList", func(t *testing.T) {
		expr, err := parser.Parse(reader("(1 ([] 'hello))"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parser.ListExpr).List))
	})
}
