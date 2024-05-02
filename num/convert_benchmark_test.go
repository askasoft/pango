package num

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func BenchmarkFtoaRegexTrailing(b *testing.B) {
	trailingZerosRegex := regexp.MustCompile(`\.?0+$`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trailingZerosRegex.ReplaceAllString("2.00000", "")
		trailingZerosRegex.ReplaceAllString("2.0000", "")
		trailingZerosRegex.ReplaceAllString("2.000", "")
		trailingZerosRegex.ReplaceAllString("2.00", "")
		trailingZerosRegex.ReplaceAllString("2.0", "")
		trailingZerosRegex.ReplaceAllString("2", "")
	}
}

func BenchmarkStripTrailingZeros(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StripTrailingZeros("2.00000")
		StripTrailingZeros("2.0000")
		StripTrailingZeros("2.000")
		StripTrailingZeros("2.00")
		StripTrailingZeros("2.0")
		StripTrailingZeros("2")
	}
}

func BenchmarkFmtF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%f", 2.03584)
	}
}

func BenchmarkStrconvF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.FormatFloat(2.03584, 'f', 6, 64)
	}
}
