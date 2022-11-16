package cpt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pandafw/pango/str"
)

const (
	SaltLength      = 16
	SecretLength    = 32
	TimestampLength = 16
	TokenLength     = SaltLength + SecretLength + TimestampLength
)

var SecretChars = str.LetterNumbers

var (
	ErrTokenLength    = errors.New("invalid token length")
	ErrTokenTimestamp = errors.New("invalid token timestamp")
)

type Token struct {
	Token     string
	Secret    string
	Salt      string
	Timestamp time.Time
}

func (t *Token) String() string {
	return t.Secret + " " + t.Timestamp.Format("2006-01-02T15:04:05Z")
}

func NewToken() *Token {
	t := &Token{
		Secret: str.RandString(SecretLength, SecretChars),
	}
	t.Refresh()
	return t
}

func ParseToken(token string) (*Token, error) {
	if len(token) != TokenLength {
		return nil, ErrTokenLength
	}

	t := &Token{Token: token}

	t.Salt = token[:SaltLength]
	t.Secret = Unsalt(token[SaltLength:SaltLength+SecretLength], t.Salt)

	s := Unsalt(token[TokenLength-TimestampLength:], t.Salt)
	ts, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return nil, ErrTokenTimestamp
	}
	t.Timestamp = time.Unix(ts, 0)
	return t, nil
}

func (t *Token) Refresh() {
	t.Salt = str.RandString(SaltLength, SecretChars)
	t.Timestamp = time.Now()
	t.Token = t.Salt + t.saltSecret() + t.saltTimestamp()
}

func (t *Token) saltSecret() string {
	return Salt(t.Secret, t.Salt)
}

func (t *Token) saltTimestamp() string {
	s := fmt.Sprintf("%016X", t.Timestamp.Unix())
	return Salt(s, t.Salt)
}

func getByte(src string, i int) byte {
	size := len(src)
	if i < 0 {
		i = (i % size) + size
	} else if i >= size {
		i %= size
	}
	return src[i]
}

func Salt(src, salt string) string {
	size := len(src)

	salted := &strings.Builder{}
	salted.Grow(size)
	for i := 0; i < size; i++ {
		x := strings.IndexByte(SecretChars, getByte(src, i))
		y := strings.IndexByte(SecretChars, getByte(salt, i))
		salted.WriteByte(getByte(SecretChars, x+y))
	}

	return salted.String()
}

func Unsalt(src, salt string) string {
	size := len(src)

	unsalted := &strings.Builder{}
	for i := 0; i < size; i++ {
		x := strings.IndexByte(SecretChars, getByte(src, i))
		y := strings.IndexByte(SecretChars, getByte(salt, i))
		unsalted.WriteByte(getByte(SecretChars, x-y))
	}
	return unsalted.String()
}
