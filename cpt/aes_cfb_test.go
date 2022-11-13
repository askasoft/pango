package cpt

import (
	"fmt"
	"testing"
)

func TestAesCFBEncryptString(t *testing.T) {
	cs := []string{
		"123456789",
		"abcdefg",
	}

	ac := NewAesCFB("1234567890123456")
	for i, c := range cs {
		o, err := ac.EncryptString(c)
		if err != nil {
			t.Fatal(i, err)
		}
		fmt.Println(i, o)

		s, err := ac.DecryptString(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if s != c {
			t.Errorf("[%d] want %q, but %q", i, c, s)
		}
	}
}

func TestAesCFBEncryptData(t *testing.T) {
	cs := [][]byte{
		[]byte("123456789"),
		[]byte("abcdefg"),
	}

	ac := NewAesCFB("1234567890123456", "0987654321654321")
	for i, c := range cs {
		o, err := ac.EncryptData(c)
		if err != nil {
			t.Fatal(i, err)
		}
		fmt.Println(i, o)

		s, err := ac.EncryptData(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if string(s) != string(c) {
			t.Errorf("[%d] want %q, but %q", i, c, s)
		}
	}
}
