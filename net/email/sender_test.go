package email

import (
	"testing"
)

func TestEncodeString(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{" 一二三四五", "=?UTF-8?q?_=E4=B8=80=E4=BA=8C=E4=B8=89=E5=9B=9B=E4=BA=94?="},
	}

	for _, c := range cs {
		a := encodeString(c.s)
		if a != c.w {
			t.Errorf("Got %v, want %v", a, c.w)
		}
	}
}
