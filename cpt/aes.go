package cpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/sha256"

	"github.com/askasoft/pango/str"
)

func NewAes128GCMCryptor(secret string) (Cryptor, error) {
	return NewAesGCMCryptor(secret, 128)
}

func NewAes192GCMCryptor(secret string) (Cryptor, error) {
	return NewAesGCMCryptor(secret, 192)
}

func NewAes256GCMCryptor(secret string) (Cryptor, error) {
	return NewAesGCMCryptor(secret, 256)
}

func NewAesGCMCryptor(secret string, bits int) (Cryptor, error) {
	k, err := hkdf.Key(sha256.New, str.UnsafeBytes(secret), nil, "", bits/8)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	g, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	return &cryptor{cipher: c, blocker: aeadBlocker{g}}, nil
}

func NewAes128CBCCryptor(secret string) (Cryptor, error) {
	return NewAesCBCCryptor(secret, 128)
}

func NewAes192CBCCryptor(secret string) (Cryptor, error) {
	return NewAesCBCCryptor(secret, 192)
}

func NewAes256CBCCryptor(secret string) (Cryptor, error) {
	return NewAesCBCCryptor(secret, 256)
}

func NewAesCBCCryptor(secret string, bits int) (Cryptor, error) {
	k, err := hkdf.Key(sha256.New, str.UnsafeBytes(secret), nil, "", bits/8)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	return &cryptor{
		cipher:  c,
		blocker: cbcBlocker{c},
		padder:  NewPkcs7Padding(c.BlockSize()),
	}, nil
}
