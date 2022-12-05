package cpt

import (
	"strings"
)

func getByte(src string, i int) byte {
	size := len(src)
	if i < 0 {
		i = (i % size) + size
	} else if i >= size {
		i %= size
	}
	return src[i]
}

func Salt(chars, salt, src string) string {
	size := len(src)

	salted := &strings.Builder{}
	salted.Grow(size)
	for i := 0; i < size; i++ {
		x := strings.IndexByte(chars, getByte(src, i))
		y := strings.IndexByte(chars, getByte(salt, i))
		salted.WriteByte(getByte(chars, x+y))
	}

	return salted.String()
}

func Unsalt(chars, salt, src string) string {
	size := len(src)

	unsalted := &strings.Builder{}
	for i := 0; i < size; i++ {
		x := strings.IndexByte(chars, getByte(src, i))
		y := strings.IndexByte(chars, getByte(salt, i))
		unsalted.WriteByte(getByte(chars, x-y))
	}
	return unsalted.String()
}
