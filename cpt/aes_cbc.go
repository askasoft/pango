package cpt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/askasoft/pango/str"
)

type AesCBC struct {
	key []byte
}

func NewAes128CBC(key string) *AesCBC {
	return NewAesCBC(key, 128)
}

func NewAes192CBC(key string) *AesCBC {
	return NewAesCBC(key, 192)
}

func NewAes256CBC(key string) *AesCBC {
	return NewAesCBC(key, 256)
}

func NewAesCBC(key string, bits int) *AesCBC {
	ac := &AesCBC{
		key: []byte(CutPadKey(key, bits/8)),
	}
	return ac
}

func (ac *AesCBC) iv() []byte {
	return ac.key[:aes.BlockSize]
}

func (ac *AesCBC) EncryptString(src string) (string, error) {
	bs, err := ac.EncryptBytes(str.UnsafeBytes(src))
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bs), nil
}

func (ac *AesCBC) DecryptString(src string) (string, error) {
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

func (ac *AesCBC) EncryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCBC: EncryptBytes: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(c, ac.iv())
	pad := Pkcs7Pad(src, cbc.BlockSize())
	des = make([]byte, len(pad))
	cbc.CryptBlocks(des, pad)

	return des, nil
}

func (ac *AesCBC) DecryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("AesCBC: DecryptBytes: Panic: %v", r)
		}
	}()

	c, err := aes.NewCipher(ac.key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(c, ac.iv())
	pad := make([]byte, len(src))
	cbc.CryptBlocks(pad, src)

	des = Pkcs7Unpad(pad)
	return des, nil
}
