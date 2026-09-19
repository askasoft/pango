package cpt

import (
	"testing"

	"github.com/askasoft/pango/ran"
)

func benchmarkAesGCMEncrypt(b *testing.B, bit int) {
	b.ResetTimer()

	c := NewAesGCMCryptor("1234567890abcde", bit)

	for range b.N {
		rs := ran.RandString(63)
		es, err := c.EncryptString(rs)
		if err != nil {
			b.Fatal(bit, err)
		}

		ds, err := c.DecryptString(es)
		if err != nil {
			b.Fatal(bit, err)
		}

		if ds != rs {
			b.Fatalf("[%d] want %q, but %q", bit, rs, ds)
		}
	}
}

func BenchmarkAes128GCMEncrypt(b *testing.B) {
	benchmarkAesGCMEncrypt(b, 128)
}

func BenchmarkAes256GCMEncrypt(b *testing.B) {
	benchmarkAesGCMEncrypt(b, 256)
}

func benchmarkAesCBCEncrypt(b *testing.B, bit int) {
	b.ResetTimer()

	c := NewAesCBCCryptor("1234567890abcde", bit)

	for range b.N {
		rs := ran.RandString(64)
		es, err := c.EncryptString(rs)
		if err != nil {
			b.Fatal(bit, err)
		}

		ds, err := c.DecryptString(es)
		if err != nil {
			b.Fatal(bit, err)
		}

		if ds != rs {
			b.Fatalf("[%d] want %q, but %q", bit, rs, ds)
		}
	}
}

func BenchmarkAes128CBCEncrypt(b *testing.B) {
	benchmarkAesCBCEncrypt(b, 128)
}

func BenchmarkAes256CBCEncrypt(b *testing.B) {
	benchmarkAesCBCEncrypt(b, 256)
}
