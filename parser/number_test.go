package parser_test

import (
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Number(suite *testing.T) {
	suite.Parallel()

	suite.Run("NumberFollowedBySymbol", func(t *testing.T) {
		expr, err := parser.Parse(reader(`-12.34 hello`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parser.NumberExpr).String())
	})

	suite.Run("SimpleInteger", func(t *testing.T) {
		expr, err := parser.Parse(reader(`12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "12", expr.(parser.NumberExpr).String())
	})

	suite.Run("Signed(-)Integer", func(t *testing.T) {
		expr, err := parser.Parse(reader(`-12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "-12", expr.(parser.NumberExpr).String())
	})

	suite.Run("Signed(+)Integer", func(t *testing.T) {
		expr, err := parser.Parse(reader(`+12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "+12", expr.(parser.NumberExpr).String())
	})

	suite.Run("SimpleFloat", func(t *testing.T) {
		expr, err := parser.Parse(reader(`+12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "+12.34", expr.(parser.NumberExpr).String())
	})

	suite.Run("SignedFloat", func(t *testing.T) {
		expr, err := parser.Parse(reader(`-12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parser.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parser.NumberExpr).String())
	})
}
