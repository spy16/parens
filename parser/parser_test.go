package parser_test

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/spy16/parens/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(suite *testing.T) {
	suite.Parallel()

	suite.Run("ReaderFailure", func(t *testing.T) {
		expr, err := parser.Parse(bufio.NewReader(readerFunc(func([]byte) (int, error) {
			return 0, errors.New("failed")
		})))
		require.Error(t, err)
		assert.Nil(t, expr)
	})

	suite.Run("UnexpectedEOF", func(t *testing.T) {
		expr, err := parser.Parse(reader(")"))
		require.Error(t, err)
		assert.Nil(t, expr)
	})
}

type readerFunc func([]byte) (int, error)

func (rf readerFunc) Read(data []byte) (int, error) {
	return rf(data)
}

func reader(s string) io.RuneScanner {
	return strings.NewReader(s)
}
