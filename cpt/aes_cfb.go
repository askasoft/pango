package cpt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/askasoft/pango/str"
)

type AesCFB struct {
	key []byte
}

func NewAes128CFB(key string) *AesCFB {
	return NewAesCFB(key, 128)
}

func NewAes192CFB(key string) *AesCFB {
	return NewAesCFB(key, 192)
}

func NewAes256CFB(key string) *AesCFB {
	return NewAesCFB(key, 256)
}

func NewAesCFB(key string, bits int) *AesCFB {
	ac := &AesCFB{
		key: []byte(CutPadKey(key, bits/8)),
	}
	return ac
}

func (ac *AesCFB) iv() []byte {
	return ac.key[:aes.BlockSize]
}

func (ac *AesCFB) EncryptString(src string) (string, error) {
	bs, err := ac.EncryptBytes(str.UnsafeBytes(src))
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

	des, err := ac.DecryptBytes(bs)
	if err != nil {
		return "", err
	}

	return str.UnsafeString(des), nil
}

func (ac *AesCFB) EncryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCFB: EncryptBytes: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(c, ac.iv())
	des = make([]byte, len(src))
	cfb.XORKeyStream(des, src)

	return des, nil
}

func (ac *AesCFB) DecryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCFB: DecryptBytes: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(c, ac.iv())
	des = make([]byte, len(src))
	cfb.XORKeyStream(des, src)

	return des, nil
}
