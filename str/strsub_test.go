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
