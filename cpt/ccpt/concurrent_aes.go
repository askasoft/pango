package ccpt

import (
	"github.com/askasoft/pango/cpt"
)

func NewAes128CBCCryptor(secret string) cpt.Cryptor {
	return NewAesCBCCryptor(secret, 128)
}

func NewAes192CBCCryptor(secret string) cpt.Cryptor {
	return NewAesCBCCryptor(secret, 192)
}

func NewAes256CBCCryptor(secret string) cpt.Cryptor {
	return NewAesCBCCryptor(secret, 256)
}

func NewAesCBCCryptor(secret string, bits int) cpt.Cryptor {
	c := &cryptor{
		fe: func() cpt.Encryptor {
			return cpt.NewAesCBCEncryptor(secret, bits)
		},
		fd: func() cpt.Decryptor {
			return cpt.NewAesCBCDecryptor(secret, bits)
		},
	}
	return c
}

func NewAes128CFBCryptor(secret string) cpt.Cryptor {
	return NewAesCFBCryptor(secret, 128)
}

func NewAes192CFBCryptor(secret string) cpt.Cryptor {
	return NewAesCFBCryptor(secret, 192)
}

func NewAes256CFBCryptor(secret string) cpt.Cryptor {
	return NewAesCFBCryptor(secret, 256)
}

func NewAesCFBCryptor(secret string, bits int) cpt.Cryptor {
	c := &cryptor{
		fe: func() cpt.Encryptor {
			return cpt.NewAesCFBEncryptor(secret, bits)
		},
		fd: func() cpt.Decryptor {
			return cpt.NewAesCFBDecryptor(secret, bits)
		},
	}
	return c
}
