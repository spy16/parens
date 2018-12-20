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
		expr, err := parens.ParseOne(bufio.NewReader(readerFunc(func([]byte) (int, error) {
			return 0, errors.New("failed")
		})))
		require.Error(t, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(")"))
		require.Error(t, err)
		assert.Nil(t, expr)
	})
}

func TestParse_Vector(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("]"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("["))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("EmptyList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("[]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 0, len(expr.(parens.VectorExpr)))
	})

	suite.Run("SimpleList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("[1 \"hello\"]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.VectorExpr)))
	})

	suite.Run("NestedList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("[1 [[] 'hello]]"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.VectorExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.VectorExpr)))
	})
}

func TestParse_Symbol(suite *testing.T) {
	suite.Parallel()

	suite.Run("AlphaSymbol", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("hello"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr(""), expr)
		assert.Equal(t, "hello", string(expr.(parens.SymbolExpr)))
	})

	suite.Run("AlphaNumericSymbol", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("hello123"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr(""), expr)
		assert.Equal(t, "hello123", string(expr.(parens.SymbolExpr)))
	})

	suite.Run("NonASCIISymbol", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("π"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr(""), expr)
		assert.Equal(t, "π", string(expr.(parens.SymbolExpr)))
	})

	suite.Run("NumerciSymbol", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("1.2.3"))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.SymbolExpr(""), expr)
		assert.Equal(t, "1.2.3", string(expr.(parens.SymbolExpr)))
	})
}

func TestParse_String(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleString", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr(""), expr)
		assert.Equal(t, "hello", string(expr.(parens.StringExpr)))
	})

	suite.Run("EscapeQuote", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello\"world"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr(""), expr)
		assert.Equal(t, "hello\"world", string(expr.(parens.StringExpr)))
	})

	suite.Run("EscapeNewline", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello\nworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr(""), expr)
		assert.Equal(t, "hello\nworld", string(expr.(parens.StringExpr)))
	})

	suite.Run("EscapeTab", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello\tworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr(""), expr)
		assert.Equal(t, "hello\tworld", string(expr.(parens.StringExpr)))
	})

	suite.Run("EscapeCR", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello\rworld"`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.StringExpr(""), expr)
		assert.Equal(t, "hello\rworld", string(expr.(parens.StringExpr)))
	})

	suite.Run("PrematureEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`"hello`))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

}

func TestParse_Number(suite *testing.T) {
	suite.Parallel()

	suite.Run("NumberFollowedBySymbol", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`-12.34 hello`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parens.NumberExpr).NumStr)
	})

	suite.Run("SimpleInteger", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "12", expr.(parens.NumberExpr).NumStr)
	})

	suite.Run("Signed(-)Integer", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`-12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12", expr.(parens.NumberExpr).NumStr)
	})

	suite.Run("Signed(+)Integer", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`+12`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "+12", expr.(parens.NumberExpr).NumStr)
	})

	suite.Run("SimpleFloat", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`+12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "+12.34", expr.(parens.NumberExpr).NumStr)
	})

	suite.Run("SignedFloat", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`-12.34`))
		assert.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.NumberExpr{}, expr)
		assert.Equal(t, "-12.34", expr.(parens.NumberExpr).NumStr)
	})
}

func TestParse_Keyword(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleKeyword", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`:hello world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.KeywordExpr(""), expr)
		assert.Equal(t, "hello", string(expr.(parens.KeywordExpr)))
	})

	suite.Run("ComplexKeyword", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`:hello/world`))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.KeywordExpr(""), expr)
		assert.Equal(t, "hello/world", string(expr.(parens.KeywordExpr)))
	})

	suite.Run("SkipsEscape", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(`:hello\nworld`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected character")
		assert.Nil(t, expr)
	})
}

func TestParse_List(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader(")"))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("("))
		require.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, expr)
	})

	suite.Run("EmptyList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("()"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 0, len(expr.(parens.ListExpr)))
	})

	suite.Run("SimpleList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("(1 \"hello\")"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.ListExpr)))
	})

	suite.Run("NestedList", func(t *testing.T) {
		expr, err := parens.ParseOne(reader("(1 ([] 'hello))"))
		require.NoError(t, err)
		require.NotNil(t, expr)
		require.IsType(t, parens.ListExpr{}, expr)
		assert.Equal(t, 2, len(expr.(parens.ListExpr)))
	})
}

type readerFunc func([]byte) (int, error)

func (rf readerFunc) Read(data []byte) (int, error) {
	return rf(data)
}

func reader(s string) io.RuneScanner {
	return strings.NewReader(s)
}
