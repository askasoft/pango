package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadLeftRune(t *testing.T) {
	assert.Equal(t, "     ", PadLeftRune("", 5, ' '))
	assert.Equal(t, "  abc", PadLeftRune("abc", 5, ' '))
	assert.Equal(t, "xxabc", PadLeftRune("abc", 5, 'x'))
	assert.Equal(t, "\uffff\uffffabc", PadLeftRune("abc", 5, '\uffff'))
	assert.Equal(t, "abc", PadLeftRune("abc", 2, ' '))
}

func TestPadLeft(t *testing.T) {
	assert.Equal(t, "     ", PadLeft("", 5, " "))
	assert.Equal(t, "-+-+abc", PadLeft("abc", 7, "-+"))
	assert.Equal(t, "-+~abc", PadLeft("abc", 6, "-+~"))
	assert.Equal(t, "-+abc", PadLeft("abc", 5, "-+~"))
	assert.Equal(t, "abc", PadLeft("abc", 2, " "))
	assert.Equal(t, "abc", PadLeft("abc", -1, " "))
	assert.Equal(t, "abc", PadLeft("abc", 5, ""))
	assert.Equal(t, "aあaあaabc", PadLeft("abc", 8, "aあ"))
}

func TestPadRightRune(t *testing.T) {
	assert.Equal(t, "     ", PadRightRune("", 5, ' '))
	assert.Equal(t, "abc  ", PadRightRune("abc", 5, ' '))
	assert.Equal(t, "abc", PadRightRune("abc", 2, ' '))
	assert.Equal(t, "abc", PadRightRune("abc", -1, ' '))
	assert.Equal(t, "abcxx", PadRightRune("abc", 5, 'x'))
}

func TestPadRight(t *testing.T) {
	assert.Equal(t, "     ", PadRight("", 5, " "))
	assert.Equal(t, "abc-+-+", PadRight("abc", 7, "-+"))
	assert.Equal(t, "abc-+~", PadRight("abc", 6, "-+~"))
	assert.Equal(t, "abc-+", PadRight("abc", 5, "-+~"))
	assert.Equal(t, "abc", PadRight("abc", 2, " "))
	assert.Equal(t, "abc", PadRight("abc", -1, " "))
	assert.Equal(t, "abc", PadRight("abc", 5, ""))
	assert.Equal(t, "abcaあaあa", PadRight("abc", 8, "aあ"))
}
