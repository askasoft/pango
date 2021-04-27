package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstrAfterByte(t *testing.T) {
	assert.Equal(t, "ot", SubstrAfterByte("foot", 'o'))
	assert.Equal(t, "bc", SubstrAfterByte("abc", 'a'))
	assert.Equal(t, "cba", SubstrAfterByte("abcba", 'b'))
	assert.Equal(t, "", SubstrAfterByte("abc", 'c'))
	assert.Equal(t, "", SubstrAfterByte("abc", 'd'))
}

func TestSubstrAfterRune(t *testing.T) {
	assert.Equal(t, "ot", SubstrAfterRune("foot", 'o'))
	assert.Equal(t, "bc", SubstrAfterRune("abc", 'a'))
	assert.Equal(t, "cba", SubstrAfterRune("abcba", 'b'))
	assert.Equal(t, "", SubstrAfterRune("abc", 'c'))
	assert.Equal(t, "", SubstrAfterRune("abc", 'd'))
}

func TestSubstrAfterLastByte(t *testing.T) {
	assert.Equal(t, "", SubstrAfterLastByte("", 'a'))
	assert.Equal(t, "", SubstrAfterLastByte("foo", 'b'))
	assert.Equal(t, "t", SubstrAfterLastByte("foot", 'o'))
	assert.Equal(t, "bc", SubstrAfterLastByte("abc", 'a'))
	assert.Equal(t, "a", SubstrAfterLastByte("abcba", 'b'))
	assert.Equal(t, "", SubstrAfterLastByte("abc", 'c'))
	assert.Equal(t, "", SubstrAfterLastByte("", 'd'))
}

func TestSubstrAfterLast(t *testing.T) {
	assert.Equal(t, "baz", SubstrAfterLast("fooXXbarXXbaz", "XX"))

	assert.Equal(t, "", SubstrAfterLast("", ""))
	assert.Equal(t, "", SubstrAfterLast("", "a"))

	assert.Equal(t, "", SubstrAfterLast("foo", "b"))
	assert.Equal(t, "t", SubstrAfterLast("foot", "o"))
	assert.Equal(t, "bc", SubstrAfterLast("abc", "a"))
	assert.Equal(t, "a", SubstrAfterLast("abcba", "b"))
	assert.Equal(t, "", SubstrAfterLast("abc", "c"))
	assert.Equal(t, "", SubstrAfterLast("", "d"))
	assert.Equal(t, "", SubstrAfterLast("abc", ""))
}

func TestSubstrBeforeByte(t *testing.T) {
	assert.Equal(t, "f", SubstrBeforeByte("foot", 'o'))
	assert.Equal(t, "", SubstrBeforeByte("abc", 'a'))
	assert.Equal(t, "a", SubstrBeforeByte("abcba", 'b'))
	assert.Equal(t, "ab", SubstrBeforeByte("abc", 'c'))
	assert.Equal(t, "", SubstrBeforeByte("abc", 'd'))
}

func TestSubstrBeforeRune(t *testing.T) {
	assert.Equal(t, "f", SubstrBeforeRune("foot", 'o'))
	assert.Equal(t, "", SubstrBeforeRune("abc", 'a'))
	assert.Equal(t, "a", SubstrBeforeRune("abcba", 'b'))
	assert.Equal(t, "ab", SubstrBeforeRune("abc", 'c'))
	assert.Equal(t, "", SubstrBeforeRune("abc", 'd'))
}

func TestSubstrBeforeLastByte(t *testing.T) {
	assert.Equal(t, "", SubstrBeforeLastByte("", 'a'))
	assert.Equal(t, "", SubstrBeforeLastByte("foo", 'b'))
	assert.Equal(t, "fo", SubstrBeforeLastByte("foot", 'o'))
	assert.Equal(t, "", SubstrBeforeLastByte("abc", 'a'))
	assert.Equal(t, "abc", SubstrBeforeLastByte("abcba", 'b'))
	assert.Equal(t, "ab", SubstrBeforeLastByte("abc", 'c'))
	assert.Equal(t, "", SubstrBeforeLastByte("", 'd'))
}

func TestSubstrBeforeLast(t *testing.T) {
	assert.Equal(t, "fooXXbar", SubstrBeforeLast("fooXXbarXXbaz", "XX"))

	assert.Equal(t, "", SubstrBeforeLast("", ""))
	assert.Equal(t, "", SubstrBeforeLast("", "a"))

	assert.Equal(t, "", SubstrBeforeLast("foo", "b"))
	assert.Equal(t, "fo", SubstrBeforeLast("foot", "o"))
	assert.Equal(t, "", SubstrBeforeLast("abc", "a"))
	assert.Equal(t, "abc", SubstrBeforeLast("abcba", "b"))
	assert.Equal(t, "ab", SubstrBeforeLast("abc", "c"))
	assert.Equal(t, "", SubstrBeforeLast("", "d"))
	assert.Equal(t, "abc", SubstrBeforeLast("abc", ""))
}
