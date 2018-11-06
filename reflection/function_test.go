package reflection_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens/reflection"
)

func add2(a, b int) int {
	return a + b
}

func addAll(vals ...float64) float64 {
	sum := 0.0
	for _, val := range vals {
		sum += val
	}
	return sum
}

var addFunc = reflect.ValueOf(add2)
var sumFunc = reflect.ValueOf(addAll)

func BenchmarkNonVariadicCall(suite *testing.B) {
	suite.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add2(1, 2)
		}
	})

	suite.Run("Reflection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflection.Call(addFunc, 1, 2)
		}
	})

	suite.Run("WithTypeConversion", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflection.Call(addFunc, []interface{}{int(1), int(2)}...)
		}
	})
}

func BenchmarkVariadicCall(suite *testing.B) {
	suite.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			addAll(1, 2)
		}
	})

	suite.Run("Reflection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflection.Call(sumFunc, 1, 2)
		}
	})

	suite.Run("WithTypeConversion", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflection.Call(sumFunc, []interface{}{int(1), int(2)}...)
		}
	})
}
