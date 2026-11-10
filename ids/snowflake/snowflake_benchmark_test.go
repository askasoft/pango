package snowflake

import (
	"testing"
)

func BenchmarkNext(b *testing.B) {
	node := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for b.Loop() {
		_ = node.NextID()
	}
}

func BenchmarkNextMaxSequence(b *testing.B) {
	node := CustomNode(DefaultEpoch, 1, 1, 21)

	b.ReportAllocs()

	b.ResetTimer()
	for b.Loop() {
		_ = node.NextID()
	}
}
