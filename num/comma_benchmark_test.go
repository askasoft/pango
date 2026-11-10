package num

import (
	"testing"
)

func BenchmarkComma(b *testing.B) {
	for b.Loop() {
		Comma(1234567890)
	}
}

func BenchmarkCommaFloat(b *testing.B) {
	for b.Loop() {
		CommaFloat(1234567890.83584, ",")
	}
}
