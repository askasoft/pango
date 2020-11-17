package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitAnyByte(t *testing.T) {
	exp := [...]string{"http", "a", "b", "c"}
	assert.Equal(t, exp[:], SplitAnyByte("http://a.b.c", ":/."), "-")
}

func TestRemoveByte(t *testing.T) {
	// RemoveByte("", *) = ""
	assert.Equal(t, "", RemoveByte("", 'a'))
	assert.Equal(t, "", RemoveByte("", 'a'))
	assert.Equal(t, "", RemoveByte("", 'a'))

	// RemoveByte("queued", 'u') = "qeed"
	assert.Equal(t, "qeed", RemoveByte("queued", 'u'))

	// RemoveByte("queued", 'z') = "queued"
	assert.Equal(t, "queued", RemoveByte("queued", 'z'))
}

func TestRemoveAnyBytes(t *testing.T) {
	// RemoveAnyByte("", *) = ""
	assert.Equal(t, "", RemoveAnyByte("", "ab"))
	assert.Equal(t, "", RemoveAnyByte("", "ab"))
	assert.Equal(t, "", RemoveAnyByte("", "ab"))

	// RemoveAnyByte("queued", 'ud') = "qee"
	assert.Equal(t, "qee", RemoveAnyByte("queued", "ud"))

	// RemoveAnyByte("queued", "z") = "queued"
	assert.Equal(t, "queued", RemoveAnyByte("queued", "z"))
}
