package reflection_test

import (
	"reflect"
	"testing"

	"github.com/spy16/parens/reflection"
)

func add(a, b int) int {
	return a + b
}

var fVal = reflect.ValueOf(add)

func BenchmarkReflectionCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflection.Call(fVal, 1, 2)
	}
}

func BenchmarkNormalCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(1, 2)
	}
}
