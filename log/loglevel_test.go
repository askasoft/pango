package log

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	cs := []struct {
		w Level
		s string
	}{
		{LevelNone, "X"},
		{LevelFatal, "F"},
		{LevelFatal, "f"},
		{LevelError, "E"},
		{LevelError, "e"},
		{LevelWarn, "W"},
		{LevelWarn, "w"},
		{LevelInfo, "I"},
		{LevelInfo, "i"},
		{LevelDebug, "D"},
		{LevelDebug, "d"},
		{LevelTrace, "T"},
		{LevelTrace, "t"},
	}

	for i, c := range cs {
		a := ParseLevel(c.s)
		if a != c.w {
			t.Errorf("[%d] ParseLevel(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
