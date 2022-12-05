package cpt

import "testing"

func TestTokenParse(t *testing.T) {
	for i := 0; i < 100; i++ {
		t1 := NewToken()

		t2, err := ParseToken(t1.token)
		if err != nil {
			t.Fatalf("[%d] ParseToken(%q)=%v", i, t1.token, err)
		}

		if t1.salt != t2.salt {
			t.Fatalf("[%d] t1.Salt = %q, t2.Salt = %q", i, t1.salt, t2.salt)
		}
		if t1.timestamp.Unix() != t2.timestamp.Unix() {
			t.Fatalf("[%d] t1.Timestamp = %d, t2.Timestamp = %d", i, t1.timestamp.Unix(), t2.timestamp.Unix())
		}
		if t1.secret != t2.secret {
			t.Fatalf("[%d] t1.Secret = %q, t2.Secret = %q", i, t1.secret, t2.secret)
		}
	}
}
