package parens

import (
	"reflect"
	"testing"
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
			reflectCall(addFunc, 1, 2)
		}
	})

	suite.Run("WithTypeConversion", func(b *testing.B) {
		// addFunc expects float64 but forcefully passing int64 which
		// triggers type-conversion
		for i := 0; i < b.N; i++ {
			reflectCall(addFunc, []interface{}{int64(1), int64(2)}...)
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
			reflectCall(sumFunc, 1, 2)
		}
	})

	suite.Run("WithTypeConversion", func(b *testing.B) {
		// sumFunc expects float64 but forcefully passing int64 which
		// triggers type-conversion
		for i := 0; i < b.N; i++ {
			reflectCall(sumFunc, []interface{}{int64(1), int64(2)}...)
		}
	})
}
