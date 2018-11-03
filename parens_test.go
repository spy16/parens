package parens_test

import (
	"errors"
	"testing"

	"github.com/spy16/parens"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute_Success(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockSExp(10, nil), nil)

	res, err := par.Execute("10")
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res, 10)
}

func TestExecute_EvalFailure(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockSExp(nil, errors.New("failed")), nil)

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func TestExecute_ParseFailure(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(nil, errors.New("failed"))

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func mockSExp(v interface{}, err error) parser.SExp {
	return sexpMock(func(scope *reflection.Scope) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}

func mockParseFn(sexp parser.SExp, err error) parens.ParseFn {
	return func(src string) (parser.SExp, error) {
		if err != nil {
			return nil, err
		}
		return sexp, nil
	}
}

type sexpMock func(scope *reflection.Scope) (interface{}, error)

func (sm sexpMock) Eval(scope *reflection.Scope) (interface{}, error) {
	return sm(scope)
}
