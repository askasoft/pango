package strutil

import (
	"regexp"
	"testing"
	"unicode"
)

func isASCIIRange(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func BenchmarkIsASCIIRange(b *testing.B) {
	str := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := isASCIIRange(str)
		if !is {
			b.Fatal("notASCII")
		}
	}
}

func BenchmarkIsASCIIIndex(b *testing.B) {
	str := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsASCII(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

var reASCII = "^[\x00-\x7F]+$"
var rxASCII = regexp.MustCompile(reASCII)

func isASCIIRegex(str string) bool {
	return rxASCII.MatchString(str)
}

func BenchmarkIsASCIIRegex(b *testing.B) {
	str := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := isASCIIRegex(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

func ascii() string {
	byt := make([]byte, unicode.MaxASCII+1)
	for i := range byt {
		byt[i] = byte(i)
	}
	return string(byt)
}
