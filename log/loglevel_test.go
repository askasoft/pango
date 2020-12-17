package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	assert.Equal(t, LevelNone, ParseLevel("X"))
	assert.Equal(t, LevelFatal, ParseLevel("F"))
	assert.Equal(t, LevelFatal, ParseLevel("f"))
	assert.Equal(t, LevelError, ParseLevel("E"))
	assert.Equal(t, LevelError, ParseLevel("e"))
	assert.Equal(t, LevelWarn, ParseLevel("W"))
	assert.Equal(t, LevelWarn, ParseLevel("w"))
	assert.Equal(t, LevelInfo, ParseLevel("I"))
	assert.Equal(t, LevelInfo, ParseLevel("i"))
	assert.Equal(t, LevelDebug, ParseLevel("D"))
	assert.Equal(t, LevelDebug, ParseLevel("d"))
	assert.Equal(t, LevelTrace, ParseLevel("T"))
	assert.Equal(t, LevelTrace, ParseLevel("t"))
}
