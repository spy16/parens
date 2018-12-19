package parser_test

import (
	"io"
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_String(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleString", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.StringExpr{}, expr)
		assert.Equal(t, "hello", expr.(parser.StringExpr).Value)
	})

	suite.Run("EscapeQuote", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello\"world"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.StringExpr{}, expr)
		assert.Equal(t, "hello\"world", expr.(parser.StringExpr).Value)
	})

	suite.Run("EscapeNewline", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello\nworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.StringExpr{}, expr)
		assert.Equal(t, "hello\nworld", expr.(parser.StringExpr).Value)
	})

	suite.Run("EscapeTab", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello\tworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.StringExpr{}, expr)
		assert.Equal(t, "hello\tworld", expr.(parser.StringExpr).Value)
	})

	suite.Run("EscapeCR", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello\rworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.StringExpr{}, expr)
		assert.Equal(t, "hello\rworld", expr.(parser.StringExpr).Value)
	})

	suite.Run("PrematureEOF", func(t *testing.T) {
		expr, err := parser.Parse(reader(`"hello`))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

}
