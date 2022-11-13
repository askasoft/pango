package cpt

import "testing"

func TestTokenSalt(t *testing.T) {
	for i := 0; i < 100; i++ {
		token := NewToken()

		salted := Salt(token.Secret, token.Salt)
		if len(token.Secret) != len(salted) {
			t.Fatalf("[%d] len(secret)=%d, len(salted)=%d", i, len(token.Secret), len(salted))
		}

		unsalted := Unsalt(salted, token.Salt)
		if token.Secret != unsalted {
			t.Fatalf("[%d] secret = %q, want %q", i, token.Secret, unsalted)
		}
	}
}

func TestTokenParse(t *testing.T) {
	for i := 0; i < 100; i++ {
		t1 := NewToken()

		t2, err := ParseToken(t1.Token)
		if err != nil {
			t.Fatalf("[%d] ParseToken(%q)=%v", i, t1.Token, err)
		}

		if t1.Secret != t2.Secret {
			t.Fatalf("[%d] t1.Secret = %q, t2.Secret = %q", i, t1.Secret, t2.Secret)
		}
		if t1.Salt != t2.Salt {
			t.Fatalf("[%d] t1.Salt = %q, t2.Salt = %q", i, t1.Salt, t2.Salt)
		}
		if t1.Timestamp.Unix() != t2.Timestamp.Unix() {
			t.Fatalf("[%d] t1.Timestamp = %d, t2.Timestamp = %d", i, t1.Timestamp.Unix(), t2.Timestamp.Unix())
		}
	}
}
