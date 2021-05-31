package vad

import (
	"regexp"
	"strings"
	"testing"
)

var (
	testReHasLowerCase = ".*[[:lower:]]"
	testReHasUpperCase = ".*[[:upper:]]"
	testRxHasLowerCase = regexp.MustCompile(testReHasLowerCase)
	testRxHasUpperCase = regexp.MustCompile(testReHasUpperCase)
)

func testLower() string {
	byt := make([]byte, 'z'-'a'+1)
	for i := 0; i < len(byt); i++ {
		byt[i] = byte(i + 'a')
	}
	return string(byt)
}

// IsLowerCase checks if the string is lowercase. Empty string is valid.
func testIsLowerCase(str string) bool {
	return str == strings.ToLower(str)
}

func BenchmarkIsLowerCase0(b *testing.B) {
	str := testLower()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsLowerCase(str)
		if !is {
			b.Fatal("notLower")
		}
	}
}

func BenchmarkIsLowerCase1(b *testing.B) {
	str := testLower()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsLowerCase(str)
		if !is {
			b.Fatal("notLower")
		}
	}
}

func testUpper() string {
	byt := make([]byte, 'Z'-'A'+1)
	for i := 0; i < len(byt); i++ {
		byt[i] = byte(i + 'A')
	}
	return string(byt)
}

// IsUpperCase checks if the string is uppercase. Empty string is valid.
func testIsUpperCase(str string) bool {
	return str == strings.ToUpper(str)
}

func BenchmarkIsUpperCase0(b *testing.B) {
	str := testUpper()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsUpperCase(str)
		if !is {
			b.Fatal("notUpper")
		}
	}
}

func BenchmarkIsUpperCase1(b *testing.B) {
	str := testUpper()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsUpperCase(str)
		if !is {
			b.Fatal("notUpper")
		}
	}
}

// HasLowerCase checks if the string contains at least 1 lowercase. Empty string is valid.
func testHasLowerCase(str string) bool {
	return testRxHasLowerCase.MatchString(str)
}

func BenchmarkHasLowerCase0(b *testing.B) {
	str := testUpper()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsLowerCase(str)
		if is {
			b.Fatal("hasLower")
		}
	}
}

func BenchmarkHasLowerCase1(b *testing.B) {
	str := testUpper()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsLowerCase(str)
		if is {
			b.Fatal("hasLower")
		}
	}
}

// HasUpperCase checks if the string contains as least 1 uppercase. Empty string is valid.
func hasUpperCase(str string) bool {
	return testRxHasUpperCase.MatchString(str)
}

func BenchmarkHasUpperCase0(b *testing.B) {
	str := testLower()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := testIsUpperCase(str)
		if is {
			b.Fatal("hasUpper")
		}
	}
}

func BenchmarkHasUpperCase1(b *testing.B) {
	str := testLower()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsUpperCase(str)
		if is {
			b.Fatal("hasUpper")
		}
	}
}
