package reflection_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spy16/parens/reflection"
)

func TestScope_Bind(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleBind", func(t *testing.T) {
		scope := reflection.NewScope(nil)
		scope.Bind("version", "1.0.0")

		val, err := scope.Value("version")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, val.RVal.Kind())
		assert.Equal(t, "1.0.0", val.RVal.String())
	})

	suite.Run("FunctionBind", func(t *testing.T) {
		scope := reflection.NewScope(nil)
		scope.Bind("print", func(msg string) { fmt.Println(msg) })

		val, err := scope.Value("print")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.Func, val.RVal.Kind())
	})

	suite.Run("OverwritingBund", func(t *testing.T) {
		scope := reflection.NewScope(nil)
		scope.Bind("print", func(msg string) { fmt.Println(msg) })
		scope.Bind("print", "now-a-string")

		val, err := scope.Value("print")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, val.RVal.Kind())
		assert.Equal(t, "now-a-string", val.RVal.String())
	})
}

func TestScope_Get(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnboundName", func(t *testing.T) {
		scope := reflection.NewScope(nil)

		val, err := scope.Value("some-unknown-name")
		require.Error(t, err)
		assert.Nil(t, val)
	})

	suite.Run("BoundOnParent", func(t *testing.T) {
		parent := reflection.NewScope(nil)
		parent.Bind("message", "hello world")

		scope := reflection.NewScope(parent)
		val, err := scope.Value("message")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, val.RVal.Kind())
		assert.Equal(t, "hello world", val.RVal.String())
	})

	suite.Run("Pointer", func(t *testing.T) {
		actualValue := "hello"

		scope := reflection.NewScope(nil)
		scope.Bind("value", &actualValue)

		val, err := scope.Value("value")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.Ptr, val.RVal.Kind())
		assert.Equal(t, &actualValue, val.RVal.Interface())
	})
}
