package num

import (
	"testing"
)

func BenchmarkParseSize(b *testing.B) {
	for b.Loop() {
		for _, s := range []string{
			"", "32", "32b", "32 B", "32k", "32.5 K", "32kb", "32 Kb",
			"32.8Mb", "32.9Gb", "32.777Tb", "32Pb", "0.3Mb", "-1",
		} {
			ParseSize(s)
			ParseSize(s)
		}
	}
}
