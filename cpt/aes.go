package cpt

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAes128CBCEncryptor(secret string) Encryptor {
	return NewAesCBCEncryptor(secret, 128)
}

func NewAes192CBCEncryptor(secret string) Encryptor {
	return NewAesCBCEncryptor(secret, 192)
}

func NewAes256CBCEncryptor(secret string) Encryptor {
	return NewAesCBCEncryptor(secret, 256)
}

func NewAesCBCEncryptor(secret string, bits int) Encryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	e := &encryptor{
		cipher:  c,
		blocker: cipher.NewCBCEncrypter(c, k[:aes.BlockSize]),
		padder:  NewPkcs7Padding(aes.BlockSize),
	}
	return e
}

func NewAes128CBCDecryptor(secret string) Decryptor {
	return NewAesCBCDecryptor(secret, 128)
}

func NewAes192CBCDecryptor(secret string) Decryptor {
	return NewAesCBCDecryptor(secret, 192)
}

func NewAes256CBCDecryptor(secret string) Decryptor {
	return NewAesCBCDecryptor(secret, 256)
}

func NewAesCBCDecryptor(secret string, bits int) Decryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	d := &decryptor{
		cipher:  c,
		blocker: cipher.NewCBCDecrypter(c, k[:aes.BlockSize]),
		padder:  NewPkcs7Padding(aes.BlockSize),
	}
	return d
}

func NewAes128CFBEncryptor(secret string) Encryptor {
	return NewAesCFBEncryptor(secret, 128)
}

func NewAes192CFBEncryptor(secret string) Encryptor {
	return NewAesCFBEncryptor(secret, 192)
}

func NewAes256CFBEncryptor(secret string) Encryptor {
	return NewAesCFBEncryptor(secret, 256)
}

func NewAesCFBEncryptor(secret string, bits int) Encryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	e := &encryptor{
		cipher:  c,
		blocker: &streamblocker{cipher.NewCFBEncrypter(c, k[:aes.BlockSize])},
		padder:  NewPkcs7Padding(aes.BlockSize),
	}
	return e
}

func NewAes128CFBDecryptor(secret string) Decryptor {
	return NewAesCFBDecryptor(secret, 128)
}

func NewAes192CFBDecryptor(secret string) Decryptor {
	return NewAesCFBDecryptor(secret, 192)
}

func NewAes256CFBDecryptor(secret string) Decryptor {
	return NewAesCFBDecryptor(secret, 256)
}

func NewAesCFBDecryptor(secret string, bits int) Decryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	d := &decryptor{
		cipher:  c,
		blocker: &streamblocker{cipher.NewCFBDecrypter(c, k[:aes.BlockSize])},
		padder:  NewPkcs7Padding(aes.BlockSize),
	}
	return d
}
