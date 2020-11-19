package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestRemoveAny(t *testing.T) {
	// RemoveAny("", *) = ""
	assert.Equal(t, "", RemoveAny("", "ab"))
	assert.Equal(t, "", RemoveAny("", "ab"))
	assert.Equal(t, "", RemoveAny("", "ab"))

	assert.Equal(t, "qee", RemoveAny("queued", "ud"))
	assert.Equal(t, "queued", RemoveAny("queued", "z"))
	assert.Equal(t, "ありとういます。", RemoveAny("ありがとうございます。", "がござ"))
}

func TestSplitAny(t *testing.T) {
	assert.Equal(t, []string{"http", "a", "b", "c"}, SplitAny("http://a.b.c", ":/."))
	assert.Equal(t, []string{"http", "あ", "い", "う"}, SplitAny("http://あ.い.う", ":/."))
	assert.Equal(t, []string{"http", "あ", "い", "う"}, SplitAny("http://あ。い。う", ":/。."))
}
