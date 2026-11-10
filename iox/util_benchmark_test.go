package iox

import (
	"io"
	"testing"

	"github.com/askasoft/pango/str"
)

type testdis struct{}

func (testdis) Write(p []byte) (int, error) {
	return len(p), nil
}

func BenchmarkWriteStringUs(b *testing.B) {
	o := &testdis{}
	s := str.Repeat("=", 1000)
	for b.Loop() {
		_, _ = WriteString(o, s)
	}
}

func BenchmarkWriteStringGo(b *testing.B) {
	o := &testdis{}
	s := str.Repeat("=", 1000)
	for b.Loop() {
		_, _ = io.WriteString(o, s)
	}
}
