package zho

import (
	"testing"
)

func TestDetectChinese(t *testing.T) {
	cs := []struct {
		w Kind
		s string
	}{
		{Hans, "真相了 “朋友不会让自己的朋友死这么多人”"},
		{Hant, "路邊見逆天美腿！"},
	}

	for i, c := range cs {
		a := DetectChinese(c.s)
		if c.w != a {
			t.Errorf("#%d DetectChinese(%q) = %v, WANT %v", i, c.s, a, c.w)
		}
	}
}
