package parens

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScope_Bind(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleBind", func(t *testing.T) {
		scope := NewScope(nil)
		scope.Bind("version", "1.0.0")

		val, err := scope.Get("version")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, reflect.TypeOf(val).Kind())
		assert.Equal(t, "1.0.0", val)
	})

	suite.Run("FunctionBind", func(t *testing.T) {
		scope := NewScope(nil)
		scope.Bind("print", func(msg string) { fmt.Println(msg) })

		val, err := scope.Get("print")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.Func, reflect.TypeOf(val).Kind())
	})

	suite.Run("OverwritingBund", func(t *testing.T) {
		scope := NewScope(nil)
		scope.Bind("print", func(msg string) { fmt.Println(msg) })
		scope.Bind("print", "now-a-string")

		val, err := scope.Get("print")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, reflect.TypeOf(val).Kind())
		assert.Equal(t, "now-a-string", val)
	})
}

func TestScope_Get(suite *testing.T) {
	suite.Parallel()

	suite.Run("UnboundName", func(t *testing.T) {
		scope := NewScope(nil)

		val, err := scope.Get("some-unknown-name")
		require.Error(t, err)
		assert.Nil(t, val)
	})

	suite.Run("BoundOnParent", func(t *testing.T) {
		parent := NewScope(nil)
		parent.Bind("message", "hello world")

		scope := NewScope(parent)
		val, err := scope.Get("message")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.String, reflect.TypeOf(val).Kind())
		assert.Equal(t, "hello world", val)
	})

	suite.Run("Pointer", func(t *testing.T) {
		actualValue := "hello"

		scope := NewScope(nil)
		scope.Bind("value", &actualValue)

		val, err := scope.Get("value")
		assert.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, reflect.Ptr, reflect.TypeOf(val).Kind())
		assert.Equal(t, &actualValue, val)
	})
}
