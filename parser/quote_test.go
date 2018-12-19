package parser_test

import (
	"io"
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Quote(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parser.Parse(reader("'"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("QuotedSymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader("'hello"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		assert.IsType(t, parser.QuoteExpr{}, expr)
		require.NotNil(t, expr.(parser.QuoteExpr).Expr)
		assert.IsType(t, parser.SymbolExpr{}, expr.(parser.QuoteExpr).Expr)
	})

	suite.Run("QuotedString", func(t *testing.T) {
		expr, err := parser.Parse(reader("'\"hello\""))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		assert.IsType(t, parser.QuoteExpr{}, expr)
		require.NotNil(t, expr.(parser.QuoteExpr).Expr)
		assert.IsType(t, parser.StringExpr{}, expr.(parser.QuoteExpr).Expr)
	})

	suite.Run("QuotedList", func(t *testing.T) {
		expr, err := parser.Parse(reader("'()"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		assert.IsType(t, parser.QuoteExpr{}, expr)
		require.NotNil(t, expr.(parser.QuoteExpr).Expr)
		assert.IsType(t, parser.ListExpr{}, expr.(parser.QuoteExpr).Expr)
	})

}
