package cpt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/str"
)

type AesCFB struct {
	key []byte
	iv  []byte
}

func NewAesCFB(key string, iv ...string) *AesCFB {
	ac := &AesCFB{key: []byte(key)}
	if len(iv) > 0 {
		ac.iv = []byte(iv[0])
	} else {
		ac.iv = ac.key
	}
	return ac
}

func (ac *AesCFB) SetKey(key string) {
	ac.key = []byte(key)
}

func (ac *AesCFB) SetIV(iv string) {
	ac.iv = []byte(iv)
}

func (ac *AesCFB) EncryptString(src string) (string, error) {
	bs, err := ac.EncryptData(str.UnsafeBytes(src))
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bs), nil
}

func (ac *AesCFB) DecryptString(src string) (string, error) {
	bs, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	des, err := ac.DecryptData(bs)
	if err != nil {
		return "", err
	}

	return bye.UnsafeString(des), nil
}

func (ac *AesCFB) EncryptData(src []byte) ([]byte, error) {
	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(c, ac.iv)
	des := make([]byte, len(src))
	cfb.XORKeyStream(des, src)

	return des, nil
}

func (ac *AesCFB) DecryptData(src []byte) ([]byte, error) {
	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(c, ac.iv)
	des := make([]byte, len(src))
	cfb.XORKeyStream(des, src)

	return des, nil
}
