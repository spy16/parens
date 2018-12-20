package parens_test

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/spy16/parens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(suite *testing.T) {
	suite.Parallel()

	suite.Run("ReaderFailure", func(t *testing.T) {
		expr, err := parens.Parse(bufio.NewReader(readerFunc(func([]byte) (int, error) {
			return 0, errors.New("failed")
		})))
		require.Error(t, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader(")"))
		require.Error(t, err)
		assert.Nil(t, expr)
	})
}

func TestParse_Vector(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader("]"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader("["))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("EmptyList", func(t *testing.T) {
		expr, err := parens.Parse(reader("[]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 0, len(expr.(parens.VectorExpr).List))
	})

	suite.Run("SimpleList", func(t *testing.T) {
		expr, err := parens.Parse(reader("[1 \"hello\"]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.VectorExpr).List))
	})

	suite.Run("NestedList", func(t *testing.T) {
		expr, err := parens.Parse(reader("[1 [[] 'hello]]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.VectorExpr).List))
	})
}

func TestParse_Symbol(suite *testing.T) {
	suite.Parallel()

	suite.Run("AlphaSymbol", func(t *testing.T) {
		expr, err := parens.Parse(reader("hello"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr{}, expr)
		assert.Equal(t, "hello", expr.(parens.SymbolExpr).String())
	})

	suite.Run("AlphaNumericSymbol", func(t *testing.T) {
		expr, err := parens.Parse(reader("hello123"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr{}, expr)
		assert.Equal(t, "hello123", expr.(parens.SymbolExpr).String())
	})

	suite.Run("NonASCIISymbol", func(t *testing.T) {
		expr, err := parens.Parse(reader("π"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr{}, expr)
		assert.Equal(t, "π", expr.(parens.SymbolExpr).String())
	})

	suite.Run("NumerciSymbol", func(t *testing.T) {
		expr, err := parens.Parse(reader("1.2.3"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr{}, expr)
		assert.Equal(t, "1.2.3", expr.(parens.SymbolExpr).String())
	})
}

func TestParse_String(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleString", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr{}, expr)
		assert.Equal(t, "hello", expr.(parens.StringExpr).Value)
	})

	suite.Run("EscapeQuote", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello\"world"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr{}, expr)
		assert.Equal(t, "hello\"world", expr.(parens.StringExpr).Value)
	})

	suite.Run("EscapeNewline", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello\nworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr{}, expr)
		assert.Equal(t, "hello\nworld", expr.(parens.StringExpr).Value)
	})

	suite.Run("EscapeTab", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello\tworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr{}, expr)
		assert.Equal(t, "hello\tworld", expr.(parens.StringExpr).Value)
	})

	suite.Run("EscapeCR", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello\rworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr{}, expr)
		assert.Equal(t, "hello\rworld", expr.(parens.StringExpr).Value)
	})

	suite.Run("PrematureEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader(`"hello`))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

}

func TestParse_Number(suite *testing.T) {
	suite.Parallel()

	suite.Run("NumberFollowedBySymbol", func(t *testing.T) {
		expr, err := parens.Parse(reader(`-12.34 hello`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parens.NumberExpr).String())
	})

	suite.Run("SimpleInteger", func(t *testing.T) {
		expr, err := parens.Parse(reader(`12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "12", expr.(parens.NumberExpr).String())
	})

	suite.Run("Signed(-)Integer", func(t *testing.T) {
		expr, err := parens.Parse(reader(`-12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12", expr.(parens.NumberExpr).String())
	})

	suite.Run("Signed(+)Integer", func(t *testing.T) {
		expr, err := parens.Parse(reader(`+12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "+12", expr.(parens.NumberExpr).String())
	})

	suite.Run("SimpleFloat", func(t *testing.T) {
		expr, err := parens.Parse(reader(`+12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "+12.34", expr.(parens.NumberExpr).String())
	})

	suite.Run("SignedFloat", func(t *testing.T) {
		expr, err := parens.Parse(reader(`-12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parens.NumberExpr).String())
	})
}

func TestParse_Keyword(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleKeyword", func(t *testing.T) {
		expr, err := parens.Parse(reader(`:hello world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.KeywordExpr{}, expr)
		assert.Equal(t, "hello", expr.(parens.KeywordExpr).String())
	})

	suite.Run("ComplexKeyword", func(t *testing.T) {
		expr, err := parens.Parse(reader(`:hello/world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.KeywordExpr{}, expr)
		assert.Equal(t, "hello/world", expr.(parens.KeywordExpr).String())
	})

	suite.Run("SkipsEscape", func(t *testing.T) {
		expr, err := parens.Parse(reader(`:hello\nworld`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected character")
		assert.Nil(t, expr)
	})
}

func TestParse_List(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader(")"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.Parse(reader("("))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("EmptyList", func(t *testing.T) {
		expr, err := parens.Parse(reader("()"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 0, len(expr.(parens.ListExpr).List))
	})

	suite.Run("SimpleList", func(t *testing.T) {
		expr, err := parens.Parse(reader("(1 \"hello\")"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.ListExpr).List))
	})

	suite.Run("NestedList", func(t *testing.T) {
		expr, err := parens.Parse(reader("(1 ([] 'hello))"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.ListExpr).List))
	})
}

type readerFunc func([]byte) (int, error)

func (rf readerFunc) Read(data []byte) (int, error) {
	return rf(data)
}

func reader(s string) io.RuneScanner {
	return strings.NewReader(s)
}
