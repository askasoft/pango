package str

import (
	crand "crypto/rand"
	mrand "math/rand"
	"testing"
	"time"
)

func BenchmarkMathGlobalRand(b *testing.B) {
	bs := make([]byte, 1000)
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		mrand.Read(bs)
	}
}

func BenchmarkMathSeedRand(b *testing.B) {
	bs := make([]byte, 1000)
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		seed := mrand.New(mrand.NewSource(time.Now().UnixNano()))
		seed.Read(bs)
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	bs := make([]byte, 1000)
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		crand.Read(bs)
	}
}

func BenchmarkRandStringByMathIntn(b *testing.B) {
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		MathRandStringIntn(100)
	}
}

func BenchmarkRandStringByMathRead(b *testing.B) {
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		MathRandStringRead(100)
	}
}

func BenchmarkRandStrinbByCryptoRead(b *testing.B) {
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		RandString(100)
	}
}

func MathRandStringRead(size int, chars ...string) string {
	cs := LetterDigitSymbols
	if len(chars) > 0 {
		cs = chars[0]
	}

	n := len(cs)

	bs := make([]byte, size)
	_, _ = mrand.Read(bs)

	for i := 0; i < size; i++ {
		bs[i] = cs[int(bs[i])%n]
	}

	return UnsafeString(bs)
}

func MathRandStringIntn(size int, chars ...string) string {
	cs := LetterDigitSymbols
	if len(chars) > 0 {
		cs = chars[0]
	}

	n := len(cs)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = cs[mrand.Intn(n)]
	}

	return UnsafeString(buf)
}
