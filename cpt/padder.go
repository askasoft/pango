package cpt

import (
	"bytes"
	"errors"

	"github.com/askasoft/pango/str"
)

// Padding interface defines functions Pad and Unpad implemented for PKCS #5 and
// PKCS #7 types of padding.
type Padding interface {
	Pad(p []byte) []byte
	Unpad(p []byte) ([]byte, error)
}

// Padder struct embeds attributes necessary for the padding calculation
// (e.g. block size). It implements the Padding interface.
type Padder int

// NewPkcs5Padding returns a PKCS5 padding type structure. The blocksize
// defaults to 8 bytes (64-bit).
// See https://tools.ietf.org/html/rfc2898 PKCS #5: Password-Based Cryptography.
// Specification Version 2.0
func NewPkcs5Padding() Padding {
	return Padder(8)
}

// NewPkcs7Padding returns a PKCS7 padding type structure. The blocksize is
// passed as a parameter.
// See https://tools.ietf.org/html/rfc2315 PKCS #7: Cryptographic Message
// Syntax Version 1.5.
// For example the block size for AES is 16 bytes (128 bits).
func NewPkcs7Padding(blockSize int) Padding {
	return Padder(blockSize)
}

func (p Padder) Pad(data []byte) []byte {
	return Pkcs7Pad(data, int(p))
}

func (p Padder) Unpad(data []byte) ([]byte, error) {
	n := len(data)
	if n == 0 {
		return nil, errors.New("padder: invalid data length")
	}

	pad := int(data[n-1])
	if pad == 0 {
		return nil, errors.New("padder: invalid last byte of padding")
	}

	if pad > n || pad > int(p) {
		return nil, errors.New("padder: invalid padding size")
	}

	for _, v := range data[n-pad : n-1] {
		if v != byte(pad) {
			return nil, errors.New("padder: invalid padding")
		}
	}

	return data[:n-pad], nil
}

// Pkcs7Pad returns the byte array passed as a parameter padded with bytes such that
// the new byte array will be an exact multiple of the expected block size.
// For example, if the expected block size is 8 bytes (e.g. PKCS #5) and that
// the initial byte array is:
//
//	[]byte{0x0A, 0x0B, 0x0C, 0x0D}
//
// the returned array will be:
//
//	[]byte{0x0A, 0x0B, 0x0C, 0x0D, 0x04, 0x04, 0x04, 0x04}
//
// The value of each octet of the padding is the size of the padding. If the
// array passed as a parameter is already an exact multiple of the block size,
// the original array will be padded with a full block.
func Pkcs7Pad(data []byte, size int) []byte {
	i := size - (len(data) % size)
	return append(data, bytes.Repeat([]byte{byte(i)}, i)...)
}

// Pkcs7Unpad removes the padding of a given byte array, according to the same rules
// as described in the Pad function. For example if the byte array passed as a
// parameter is:
//
//	[]byte{0x0A, 0x0B, 0x0C, 0x0D, 0x04, 0x04, 0x04, 0x04}
//
// the returned array will be:
//
//	[]byte{0x0A, 0x0B, 0x0C, 0x0D}
func Pkcs7Unpad(data []byte) []byte {
	n := len(data)
	p := int(data[n-1])
	return data[:n-p]
}

// CutPadKey cut key if key's length is greater than 'size',
// or pad key with space if key's length is smaller than 'size'.
func CutPadKey(key string, size int) string {
	if len(key) > size {
		key = key[:size]
	} else if len(key) < size {
		key = str.PadRight(key, size, " ")
	}
	return key
}
