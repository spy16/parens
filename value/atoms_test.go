package value

import (
	"testing"
)

func BenchmarkString_HashCode(b *testing.B) {
	s := &String{Value: "hello"}
	var hash int64
	for i := 0; i < b.N; i++ {
		hash = s.HashCode()
	}
	_ = dummy(hash)
}

func BenchmarkInt64_HashCode(b *testing.B) {
	i64 := Int64(83797979)
	var hash int64
	for i := 0; i < b.N; i++ {
		hash = i64.HashCode()
	}
	_ = dummy(hash)
}

func BenchmarkFloat64_HashCode(b *testing.B) {
	f64 := Float64(1.4324234)
	var hash int64
	for i := 0; i < b.N; i++ {
		hash = f64.HashCode()
	}
	_ = dummy(hash)
}

func dummy(v int64) int64 { return v }
