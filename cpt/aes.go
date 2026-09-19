package cpt

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAes128GCMCryptor(secret string) Cryptor {
	return NewAesGCMCryptor(secret, 128)
}

func NewAes192GCMCryptor(secret string) Cryptor {
	return NewAesGCMCryptor(secret, 192)
}

func NewAes256GCMCryptor(secret string) Cryptor {
	return NewAesGCMCryptor(secret, 256)
}

func NewAesGCMCryptor(secret string, bits int) Cryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	g, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	return &cryptor{
		cipher:  c,
		blocker: aeadBlocker{g},
	}
}

func NewAes128CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 128)
}

func NewAes192CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 192)
}

func NewAes256CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 256)
}

func NewAesCBCCryptor(secret string, bits int) Cryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	return &cryptor{
		cipher:  c,
		blocker: cbcBlocker{c},
		padder:  NewPkcs7Padding(c.BlockSize()),
	}
}
