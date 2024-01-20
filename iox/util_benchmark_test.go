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
	for i := 0; i < b.N; i++ {
		_, _ = WriteString(o, s)
	}
}

func BenchmarkWriteStringGo(b *testing.B) {
	o := &testdis{}
	s := str.Repeat("=", 1000)
	for i := 0; i < b.N; i++ {
		_, _ = io.WriteString(o, s)
	}
}
