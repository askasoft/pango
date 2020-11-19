package str

import (
	"regexp"
	"testing"
	"unicode"
)

func testIsASCIIRange(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func BenchmarkIsASCIIRange(b *testing.B) {
	str := testAscii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsASCIIRange(str)
		if !is {
			b.Fatal("notASCII")
		}
	}
}

func BenchmarkIsASCIIIndex(b *testing.B) {
	str := testAscii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsASCII(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

var testReASCII = "^[\x00-\x7F]+$"
var testRxASCII = regexp.MustCompile(testReASCII)

func testIsASCIIRegex(str string) bool {
	return testRxASCII.MatchString(str)
}

func BenchmarkIsASCIIRegex(b *testing.B) {
	str := testAscii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsASCIIRegex(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

func testAscii() string {
	byt := make([]byte, unicode.MaxASCII+1)
	for i := range byt {
		byt[i] = byte(i)
	}
	return string(byt)
}
