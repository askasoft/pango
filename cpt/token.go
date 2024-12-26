package cpt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/askasoft/pango/ran"
	"github.com/askasoft/pango/str"
)

const (
	SecretChars = str.Base64URL
)

var (
	ErrTokenLength    = errors.New("invalid token length")
	ErrTokenTimestamp = errors.New("invalid token timestamp")
)

var timestampLength = 16

var tokener = NewTokener(8, 16)

// RandomToken create a random `randLenth` secret token.
func RandomToken() *Token {
	return tokener.RandomToken()
}

// SecretToken create a token with secret.
func SecretToken(secret string) *Token {
	return tokener.SecretToken(secret)
}

func ParseToken(token string) (t *Token, err error) {
	return tokener.ParseToken(token)
}

type Tokener struct {
	saltLength int
	randLength int
}

// NewTokener create a tokener.
func NewTokener(saltLength, randLength int) *Tokener {
	return &Tokener{
		saltLength: saltLength,
		randLength: randLength,
	}
}

// RandomToken create a random `randLenth` secret token.
func (tr *Tokener) RandomToken() *Token {
	return tr.newToken(ran.RandString(tr.randLength, SecretChars))
}

// SecretToken create a token with secret.
// If secret is empty, a random `randLenth` secret will be used.
func (tr *Tokener) SecretToken(secret string) *Token {
	return tr.newToken(secret)
}

func (tr *Tokener) newToken(secret string) *Token {
	t := &Token{
		salt:   str.RepeatByte(' ', tr.saltLength),
		secret: secret,
	}

	t.Refresh()
	return t
}

func (tr *Tokener) ParseToken(token string) (*Token, error) {
	if len(token) < tr.saltLength+timestampLength {
		return nil, ErrTokenLength
	}

	t := &Token{
		token: token,
		salt:  token[:tr.saltLength],
	}

	ts := t.token[tr.saltLength : tr.saltLength+timestampLength]
	ts = Unsalt(SecretChars, t.salt, ts)
	tp, err := strconv.ParseInt(ts, 16, 64)
	if err != nil {
		return nil, ErrTokenTimestamp
	}
	t.timestamp = time.Unix(tp, 0)

	t.secret = t.token[tr.saltLength+timestampLength:]
	t.secret = Unsalt(SecretChars, t.salt, t.secret)
	return t, nil
}

type Token struct {
	salt      string
	timestamp time.Time
	secret    string
	token     string
}

func (t *Token) Salt() string {
	return t.salt
}

func (t *Token) Timestamp() time.Time {
	return t.timestamp
}

func (t *Token) Secret() string {
	return t.secret
}

func (t *Token) Token() string {
	return t.token
}

func (t *Token) String() string {
	return t.timestamp.Format("2006-01-02T15:04:05Z") + " " + t.secret
}

func (t *Token) Refresh() {
	t.salt = ran.RandString(len(t.salt), SecretChars)

	t.timestamp = time.Now()
	ts := fmt.Sprintf("%016x", t.timestamp.Unix())
	st := Salt(SecretChars, t.salt, ts)
	ss := Salt(SecretChars, t.salt, t.secret)
	t.token = t.salt + st + ss
}
