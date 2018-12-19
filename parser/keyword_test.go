package parser_test

import (
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Keyword(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleKeyword", func(t *testing.T) {
		expr, err := parser.Parse(reader(`:hello world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.KeywordExpr{}, expr)
		assert.Equal(t, "hello", expr.(parser.KeywordExpr).String())
	})

	suite.Run("ComplexKeyword", func(t *testing.T) {
		expr, err := parser.Parse(reader(`:hello/world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.KeywordExpr{}, expr)
		assert.Equal(t, "hello/world", expr.(parser.KeywordExpr).String())
	})

	suite.Run("SkipsEscape", func(t *testing.T) {
		expr, err := parser.Parse(reader(`:hello\nworld`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected character")
		assert.Nil(t, expr)
	})
}
