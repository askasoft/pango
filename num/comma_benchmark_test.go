package num

import (
	"testing"
)

func BenchmarkComma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Comma(1234567890)
	}
}

func BenchmarkCommaFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CommaFloat(1234567890.83584)
	}
}
