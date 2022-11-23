package cpt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/pandafw/pango/str"
)

type AesCBC struct {
	key []byte
	iv  []byte
}

func NewAesCBC(key string, iv ...string) *AesCBC {
	ac := &AesCBC{key: []byte(key)}
	if len(iv) > 0 {
		ac.iv = []byte(iv[0])
	} else {
		ac.iv = ac.key
	}
	return ac
}

func (ac *AesCBC) SetKey(key string) {
	ac.key = []byte(key)
}

func (ac *AesCBC) SetIV(iv string) {
	ac.iv = []byte(iv)
}

func (ac *AesCBC) EncryptString(src string) (string, error) {
	bs, err := ac.EncryptData(str.UnsafeBytes(src))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

func (ac *AesCBC) DecryptString(src string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	des, err := ac.DecryptData(bs)
	if err != nil {
		return "", err
	}

	return string(des), nil
}

func (ac *AesCBC) EncryptData(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCBC: EncryptData: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher([]byte(ac.key))
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(c, ac.iv)
	pad := Pkcs7Pad(src, cbc.BlockSize())
	des = make([]byte, len(pad))
	cbc.CryptBlocks(des, pad)

	return des, nil
}

func (ac *AesCBC) DecryptData(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCBC: DecryptData: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher([]byte(ac.key))
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(c, ac.iv)
	pad := make([]byte, len(src))
	cbc.CryptBlocks(pad, src)

	des = Pkcs7Unpad(pad)
	return des, nil
}
