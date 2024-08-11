package cpt

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/askasoft/pango/str"
)

type Encryptor interface {
	EncryptString(string) (string, error)
	EncryptBytes([]byte) ([]byte, error)
}

type Decryptor interface {
	DecryptString(string) (string, error)
	DecryptBytes([]byte) ([]byte, error)
}

type Cryptor interface {
	Encryptor
	Decryptor
}

type encryptor struct {
	cipher  cipher.Block
	blocker cipher.BlockMode
	padder  Padding
}

func (e *encryptor) EncryptString(src string) (string, error) {
	bs, err := e.EncryptBytes(str.UnsafeBytes(src))
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bs), nil
}

func (e *encryptor) EncryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("encryptor.EncryptBytes: Panic: %v", r)
		}
	}()

	if e.blocker.BlockSize() > 0 {
		src = e.padder.Pad(src)
	}

	des = make([]byte, len(src))
	e.blocker.CryptBlocks(des, src)
	return
}

type decryptor struct {
	cipher  cipher.Block
	blocker cipher.BlockMode
	padder  Padding
}

func (d *decryptor) DecryptString(src string) (string, error) {
	bs, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	des, err := d.DecryptBytes(bs)
	if err != nil {
		return "", err
	}

	return str.UnsafeString(des), nil
}

func (d *decryptor) DecryptBytes(src []byte) (des []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			des = nil
			err = fmt.Errorf("decryptor.DecryptBytes: Panic: %v", r)
		}
	}()

	des = make([]byte, len(src))
	d.blocker.CryptBlocks(des, src)

	if d.blocker.BlockSize() > 0 {
		des, err = d.padder.Unpad(des)
	}
	return
}

type streamblocker struct {
	stream cipher.Stream
}

// BlockSize returns the mode's block size.
func (sb *streamblocker) BlockSize() int {
	return 0
}

func (sb *streamblocker) CryptBlocks(dst, src []byte) {
	sb.stream.XORKeyStream(dst, src)
}
