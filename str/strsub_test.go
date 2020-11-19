package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringAfterByte(t *testing.T) {
	assert.Equal(t, "ot", StringAfterByte("foot", 'o'))
	assert.Equal(t, "bc", StringAfterByte("abc", 'a'))
	assert.Equal(t, "cba", StringAfterByte("abcba", 'b'))
	assert.Equal(t, "", StringAfterByte("abc", 'c'))
	assert.Equal(t, "", StringAfterByte("abc", 'd'))
}

func TestStringAfterRune(t *testing.T) {
	assert.Equal(t, "ot", StringAfterRune("foot", 'o'))
	assert.Equal(t, "bc", StringAfterRune("abc", 'a'))
	assert.Equal(t, "cba", StringAfterRune("abcba", 'b'))
	assert.Equal(t, "", StringAfterRune("abc", 'c'))
	assert.Equal(t, "", StringAfterRune("abc", 'd'))
}

func TestStringAfterLastByte(t *testing.T) {
	assert.Equal(t, "", StringAfterLastByte("", 'a'))
	assert.Equal(t, "", StringAfterLastByte("foo", 'b'))
	assert.Equal(t, "t", StringAfterLastByte("foot", 'o'))
	assert.Equal(t, "bc", StringAfterLastByte("abc", 'a'))
	assert.Equal(t, "a", StringAfterLastByte("abcba", 'b'))
	assert.Equal(t, "", StringAfterLastByte("abc", 'c'))
	assert.Equal(t, "", StringAfterLastByte("", 'd'))
}

func TestStringAfterLast(t *testing.T) {
	assert.Equal(t, "baz", StringAfterLast("fooXXbarXXbaz", "XX"))

	assert.Equal(t, "", StringAfterLast("", ""))
	assert.Equal(t, "", StringAfterLast("", "a"))

	assert.Equal(t, "", StringAfterLast("foo", "b"))
	assert.Equal(t, "t", StringAfterLast("foot", "o"))
	assert.Equal(t, "bc", StringAfterLast("abc", "a"))
	assert.Equal(t, "a", StringAfterLast("abcba", "b"))
	assert.Equal(t, "", StringAfterLast("abc", "c"))
	assert.Equal(t, "", StringAfterLast("", "d"))
	assert.Equal(t, "", StringAfterLast("abc", ""))
}
