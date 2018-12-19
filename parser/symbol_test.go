package parser_test

import (
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Symbol(suite *testing.T) {
	suite.Parallel()

	suite.Run("AlphaSymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader("hello"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.SymbolExpr{}, expr)
		assert.Equal(t, "hello", expr.(parser.SymbolExpr).String())
	})

	suite.Run("AlphaNumericSymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader("hello123"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.SymbolExpr{}, expr)
		assert.Equal(t, "hello123", expr.(parser.SymbolExpr).String())
	})

	suite.Run("NonASCIISymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader("π"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.SymbolExpr{}, expr)
		assert.Equal(t, "π", expr.(parser.SymbolExpr).String())
	})

	suite.Run("NumerciSymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader("1.2.3"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.SymbolExpr{}, expr)
		assert.Equal(t, "1.2.3", expr.(parser.SymbolExpr).String())
	})
}
