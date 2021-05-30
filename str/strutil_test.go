package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartsWith(t *testing.T) {
	assert.True(t, StartsWith("", ""))
	assert.True(t, StartsWith("foobar", ""))
	assert.False(t, StartsWith("", "f"))

	assert.True(t, StartsWith("foobar", "foo"))
	assert.True(t, StartsWith("あいうえお", "あ"))

	assert.False(t, StartsWith("f", "oo"))
	assert.False(t, StartsWith("あ", "あいうえお"))
	assert.False(t, StartsWith("foobar", "oo"))
	assert.False(t, StartsWith("あいうえお", "い"))
}

func TestEndsWith(t *testing.T) {
	assert.True(t, EndsWith("", ""))
	assert.True(t, EndsWith("foobar", ""))
	assert.False(t, EndsWith("", "f"))

	assert.True(t, EndsWith("foobar", "bar"))
	assert.True(t, EndsWith("あいうえお", "えお"))

	assert.False(t, EndsWith("f", "oo"))
	assert.False(t, EndsWith("あ", "あいうえお"))
	assert.False(t, EndsWith("foobar", "oo"))
	assert.False(t, EndsWith("あいうえお", "い"))
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

func TestRemoveAny(t *testing.T) {
	// RemoveAny("", *) = ""
	assert.Equal(t, "", RemoveAny("", "ab"))
	assert.Equal(t, "", RemoveAny("", "ab"))
	assert.Equal(t, "", RemoveAny("", "ab"))

	assert.Equal(t, "qee", RemoveAny("queued", "ud"))
	assert.Equal(t, "queued", RemoveAny("queued", "z"))
	assert.Equal(t, "ありとういます。", RemoveAny("ありがとうございます。", "がござ"))
}

func TestRemoveAnyByte(t *testing.T) {
	// RemoveAnyByte("", *) = ""
	assert.Equal(t, "", RemoveAnyByte("", "ab"))
	assert.Equal(t, "", RemoveAnyByte("", "ab"))
	assert.Equal(t, "", RemoveAnyByte("", "ab"))

	assert.Equal(t, "qee", RemoveAnyByte("queued", "ud"))
	assert.Equal(t, "queued", RemoveAnyByte("queued", "z"))
}

func TestSplitAny(t *testing.T) {
	assert.Equal(t, []string{""}, SplitAny("", "c"))
	assert.Equal(t, []string{""}, SplitAny("", ".c"))
	assert.Equal(t, []string{"http://a", "b-", ""}, SplitAny("http://a.b-c", ".c"))
	assert.Equal(t, []string{"http", "", "", "a", "b", "c"}, SplitAny("http://a.b.c", ":/."))
	assert.Equal(t, []string{"http", "", "", "あ", "い", "う"}, SplitAny("http://あ.い.う", ":/."))
	assert.Equal(t, []string{"http", "", "", "あ", "い", "う"}, SplitAny("http://あ。い。う", ":/。."))
}

func TestFieldsAny(t *testing.T) {
	assert.Equal(t, []string{}, FieldsAny("", "c"))
	assert.Equal(t, []string{}, FieldsAny("", ".c"))
	assert.Equal(t, []string{"http://a", "b"}, FieldsAny("http://a.b.c", ".c"))
	assert.Equal(t, []string{"http", "a", "b", "c"}, FieldsAny("http://a.b.c", ":/."))
	assert.Equal(t, []string{"http", "あ", "い", "う"}, FieldsAny("http://あ.い.う", ":/."))
	assert.Equal(t, []string{"http", "あ", "い", "う"}, FieldsAny("http://あ。い。う", ":/。."))
}
